# Investigation plan: remaining slow-SQL queries not covered by CN-5272

Related work: [CN-5272](https://thd.atlassian.net/browse/CN-5272) —
[docs/CN-5272/kind-casing-and-index-fix.md](CN-5272/kind-casing-and-index-fix.md)

## Where this data came from

Pulled `SLOW SQL >= 200ms` log lines (GORM's slow-query logger) from all 5
`spin-go-clouddriver` pods in the `pr-te-cd-tools` / `us-central1` (production)
cluster, `spinnaker` namespace, over a 24h window:

```bash
kubectl --context gke_pr-te-cd-tools_us-central1_pr-us-central1 -n spinnaker \
  logs <pod> --since=24h | grep -A1 "SLOW SQL"
```

6,936 slow-SQL log lines total. Grouped by source location (`internal/sql/client.go:<line>`,
line numbers per the committed `HEAD` version), ignoring the literal values bound
into each query:

| Count | Location | Function | Covered by CN-5272? |
|---|---|---|---|
| 6,103 | `client.go:367` | `ListKubernetesClustersByFields` | **Yes** |
| 53 | `client.go:348` | `ListKubernetesClustersByApplication` | **Yes** |
| 207 | `client.go:188` | `GetKubernetesProvider` | No — this doc |
| 165 | `client.go:264` | `GetKubernetesProviderAndPermissions` | No — this doc |
| 103 | `client.go:463` | `ListKubernetesProvidersAndPermissions` | No — this doc |
| 86 | `client.go:137` | `CreateKubernetesResource` | No — this doc |
| 170 | `client.go:595` | `ListKubernetesAccountsBySpinnakerApp` | No — this doc |
| 25 | `client.go:384` | `ListKubernetesProviders` | No — this doc |
| 24 | `client.go:562` | `ListKubernetesResourcesByTaskID` | No — this doc |

The two CN-5272 queries account for ~89% of all slow-SQL volume. This document
covers the remaining seven query shapes — real, unaddressed problems in their
own right, but a distinct investigation from CN-5272.

## Duration profile per query (same 24h sample)

| Location | n | min | p50 | avg | p95 | max |
|---|---|---|---|---|---|---|
| `client.go:188` | 207 | 202.8ms | 365.6ms | 859.0ms | 3220.6ms | 13618.9ms |
| `client.go:264` | 165 | 200.7ms | 346.0ms | 505.5ms | 1069.8ms | 5150.6ms |
| `client.go:463` | 103 | 201.3ms | 406.6ms | 706.1ms | 2939.2ms | 10166.7ms |
| `client.go:137` | 86 | 202.3ms | 344.4ms | 469.7ms | 982.3ms | 4089.8ms |
| `client.go:595` | 170 | 201.9ms | 386.1ms | 564.6ms | 1124.5ms | 6909.3ms |
| `client.go:384` | 25 | 202.0ms | 262.6ms | 756.5ms | 2082.5ms | 5094.1ms |
| `client.go:562` | 24 | 214.4ms | 356.6ms | 649.2ms | 2547.1ms | 3718.2ms |

**The shape of this data matters as much as the averages.** Every one of these
seven queries has a `p50` in the 250-400ms range but a `max` an order of
magnitude higher (5-13 seconds). A query with a genuinely bad plan (e.g. a full
table scan from a missing index) is typically **consistently** slow — its
duration wouldn't vary 30-50x between p50 and max. This min/p50-vs-max spread
is a strong signal of **intermittent contention** (lock waits, I/O queuing, CPU
starvation) rather than seven separate missing-index bugs. The obvious
candidate source of that contention: `client.go:367`, which ran 6,103 times in
the same 24h window, most of those as multi-second full-table scans against the
same database. **This hypothesis should be treated as the leading one for all
seven queries below**, and re-checked directly (see "Cross-cutting step 0")
before investing time in per-query fixes that may turn out to be unnecessary.

---

## Cross-cutting step 0 (do this first, before any per-query work below)

Re-pull this same log analysis **after** the CN-5272 fix (query rewrite +
index) has been deployed and had a few days of production traffic. If the
`p95`/`max` figures for some or all of these seven queries drop substantially
without any change to their own code, that confirms the contention hypothesis
for those queries and closes them out (or substantially re-scopes them) without
further work. Any query whose `p95`/`max` *doesn't* improve after CN-5272 ships
is the one genuinely worth the deep-dive investigation below — that's the
signal for prioritizing which of the seven to actually dig into first.

---

## `GetKubernetesProvider` — `client.go:188`

```sql
SELECT a.name, a.host, a.ca_data, a.bearer_token, a.token_provider, a.namespace as legacy_namespace, b.namespace
FROM kubernetes_providers a
LEFT JOIN kubernetes_providers_namespaces b ON a.name = b.account_name
WHERE a.name = ?
```

- **Highest volume (207) and highest max (13.6s) of the unaddressed group.**
- `kubernetes_providers.name` is the primary key (`internal/kubernetes/provider.go`) —
  the `WHERE a.name = ?` side should be a PK point lookup.
- `kubernetes_providers_namespaces` has only a composite unique index
  `(account_name, namespace)` today (no standalone `account_name` index in the
  currently-deployed code — that's an *uncommitted* local change from earlier
  in this work, not yet in production). The join should still be able to use
  the composite index's leftmost prefix (`account_name`), but confirm this,
  don't assume it.
- **Investigation steps:**
  1. Do cross-cutting Step 0 first.
  2. `EXPLAIN` this exact query shape (with a real account name) directly
     against prod — confirm whether `b` is actually using the composite index's
     prefix or scanning `kubernetes_providers_namespaces` outright.
  3. Check `kubernetes_providers_namespaces` row count and whether any single
     account has an unusually large number of namespace rows (a many-namespace
     account could make this join disproportionately expensive for that one
     account, explaining why some calls are fast and some are very slow).
  4. Check whether this endpoint is called per-request from a hot path (e.g.
     once per pipeline execution/account resolution) — if call *volume* is the
     real issue rather than plan quality, caching the provider lookup
     (already-fetched providers are largely static, config-file-driven) may be
     the more impactful fix than a query change.

## `GetKubernetesProviderAndPermissions` — `client.go:264`

```sql
SELECT a.name, a.host, a.ca_data, a.token_provider, a.namespace as legacy_namespace, d.namespace, b.read_group, c.write_group
FROM kubernetes_providers a
LEFT JOIN provider_read_permissions b ON a.name = b.account_name
LEFT JOIN provider_write_permissions c ON a.name = c.account_name
LEFT JOIN kubernetes_providers_namespaces d ON a.name = d.account_name
WHERE a.name = ?
```

- Same shape as `GetKubernetesProvider` plus two more joins.
- `provider_read_permissions` and `provider_write_permissions` both already
  have a standalone `account_name_idx` (`pkg/permissions.go`) — these two joins
  are less suspect than the `kubernetes_providers_namespaces` join, which again
  only has the composite unique index.
- **Investigation steps:**
  1. Do cross-cutting Step 0 first.
  2. `EXPLAIN` this query shape against prod; compare the join order MySQL
     picks for all three joined tables.
  3. Check whether any single account has a disproportionately large number of
     read/write permission group rows or namespaces — same "outlier account"
     hypothesis as `GetKubernetesProvider`.
  4. This function and `GetKubernetesProvider` are near-duplicates of each
     other (same base joins, one adds permissions) — worth asking whether both
     are needed on hot paths or whether callers could share one cached result.

## `ListKubernetesProvidersAndPermissions` — `client.go:463`

```sql
SELECT a.name, a.host, a.ca_data, a.token_provider, a.namespace as legacy_namespace, d.namespace, b.read_group, c.write_group
FROM kubernetes_providers a
LEFT JOIN provider_read_permissions b ON a.name = b.account_name
LEFT JOIN provider_write_permissions c ON a.name = c.account_name
LEFT JOIN kubernetes_providers_namespaces d ON a.name = d.account_name
```

- Same 3-way join as `GetKubernetesProviderAndPermissions`, but **unfiltered** —
  every call joins across the entire providers/permissions/namespaces table set
  at once, for every account THD has onboarded to this clouddriver instance.
- This is structurally the most expensive of the group by design (no `WHERE`
  clause at all) — the question isn't "why is this slow" so much as "how often
  is this called, and does it need to be." Second-highest max (10.2s) despite
  lowest median call frequency among the group.
- **Investigation steps:**
  1. Do cross-cutting Step 0 first, but treat this one skeptically — an
     unfiltered 3-way join across every account is a plausible source of
     genuine (not just contention-driven) slowness regardless of what else is
     happening on the box.
  2. Find every caller of `ListKubernetesProvidersAndPermissions` — grep
     `internal/api` for its usage — and determine whether it's called on a
     user-facing hot path (e.g. per-request) versus something that could be
     cached/refreshed on an interval instead (this data — which accounts exist
     and their permissions — changes rarely compared to how often it's likely
     queried).
  3. Get the actual row counts of `kubernetes_providers`,
     `provider_read_permissions`, `provider_write_permissions`, and
     `kubernetes_providers_namespaces` — an unfiltered join's cost scales with
     all four, so this determines whether the fix is "add caching" vs.
     "this table has grown and the join itself needs restructuring."

## `ListKubernetesProviders` — `client.go:384`

```sql
SELECT a.name, a.host, a.ca_data, a.token_provider, a.namespace as legacy_namespace, b.namespace
FROM kubernetes_providers a
LEFT JOIN kubernetes_providers_namespaces b ON a.name = b.account_name
```

- The unfiltered counterpart to `GetKubernetesProvider`, same relationship as
  `ListKubernetesProvidersAndPermissions` is to `GetKubernetesProviderAndPermissions`.
  Lowest volume (25) of the group but still hits a 5s max.
- **Investigation steps:** same as `ListKubernetesProvidersAndPermissions`
  above — find callers, determine if this can be cached/interval-refreshed
  instead of queried per-call, and get real row counts for the two joined
  tables.

## `ListKubernetesAccountsBySpinnakerApp` — `client.go:595`

```sql
SELECT account_name FROM kubernetes_resources WHERE spinnaker_app = ? GROUP BY account_name
```

- **This is the query CN-5272's Decision 2 explicitly deferred an index for** —
  `idx_kubernetes_resources_app_covering (spinnaker_app, account_name, cluster, kind)`
  was considered and intentionally not added without slow-query-log evidence.
  **170 occurrences over 24h, with a 6.9s max, is exactly that evidence.**
- Right now this query has no index on `spinnaker_app` at all — it's filtering
  on a column with no dedicated index, on the same large, high-churn
  `kubernetes_resources` table.
- **Investigation steps:**
  1. Do cross-cutting Step 0 first — this table is the same one query
     `client.go:367` hammers, so contention is plausible here too, but this
     query's lack of a `spinnaker_app` index is a real, independent
     candidate cause, unlike the provider/permission queries above.
  2. `EXPLAIN` this query shape against prod to confirm it's doing a full
     table scan on `kubernetes_resources` (630k+ rows as of the CN-5272 work).
  3. If Step 0 doesn't fully explain this away, this is the strongest
     candidate in the whole unaddressed group for a **direct, low-risk index
     fix**: revisit CN-5272 Decision 2 and consider shipping
     `idx_kubernetes_resources_app_covering` via the same manual-DDL runbook
     pattern already built for `idx_kubernetes_resources_kind_covering`
     ([docs/CN-5272/kind-covering-index-rollout-runbook.md](CN-5272/kind-covering-index-rollout-runbook.md)).

## `ListKubernetesResourcesByTaskID` — `client.go:562`

```sql
SELECT account_name, api_group, kind, name, artifact_name, namespace, resource, task_type, version
FROM kubernetes_resources WHERE task_id = ?
```

- Lowest volume (24) of the group, but a 3.7s max is still surprising:
  `task_id_idx` already exists on this exact column — this should be a fast,
  indexed point lookup essentially every time.
- **Investigation steps:**
  1. Do cross-cutting Step 0 first — this is the query shape most likely to be
     *purely* contention-driven, since there's no obvious missing-index
     explanation and the index it needs already exists.
  2. If it's still slow post-CN-5272, `EXPLAIN` it to confirm `task_id_idx` is
     actually being used and check whether any single `task_id` has an
     unusually large number of associated resource rows (a task touching many
     resources could make this legitimately slower for that one task, same
     "outlier" pattern as the provider queries).

## `CreateKubernetesResource` — `client.go:137`

```sql
INSERT INTO kubernetes_resources (...) VALUES (...)
```

- A single-row `INSERT` taking 200ms-4s is not a query-plan problem — inserts
  don't have an "index the WHERE clause" fix. Slowness here is inherently about
  write-path cost: index maintenance overhead, lock waits from concurrent
  writes/reads on the same table, or general I/O contention.
- **Directly relevant to CN-5272's own tradeoffs:** every index on this table
  (there are 3 today: `account_name_kind_name_spinnaker_app_idx`, `kind_idx`,
  `task_id_idx`; a 4th, `idx_kubernetes_resources_kind_covering`, is pending
  rollout) adds insert-time cost. This query is the direct, measurable cost
  side of that tradeoff — track it specifically after the new index ships (this
  was already flagged as a monitoring item in the rollout runbook's Step 11).
- **Investigation steps:**
  1. Do cross-cutting Step 0 first — if `client.go:367`'s scans are locking or
     saturating I/O on `kubernetes_resources`, concurrent inserts into the same
     table are a very plausible casualty. This is the single most likely
     "resolves mostly on its own" candidate in the whole list.
  2. If still slow post-CN-5272: check Cloud SQL IOPS/CPU during insert spikes
     specifically (not just during the old `367` scans) to see if write
     saturation is a standing issue independent of the read-side fix.
  3. Watch this metric specifically once `idx_kubernetes_resources_kind_covering`
     is applied (per the rollout runbook) — if insert latency measurably
     worsens post-rollout, that's the write-amplification cost predicted in
     CN-5272 Decision 2 materializing, and is the strongest argument against
     ever adding the deferred second index without very strong justification.

---

## Priority ranking for follow-up (post cross-cutting Step 0)

1. **`ListKubernetesAccountsBySpinnakerApp` (595)** — has a concrete, already-known
   candidate fix (the deferred index) and the clearest non-contention
   explanation (no index on the filtered column at all).
2. **`ListKubernetesProvidersAndPermissions` (463) / `ListKubernetesProviders` (384)** —
   unfiltered joins across whole tables; likely a caching/call-pattern fix
   rather than an indexing fix. Worth checking caller frequency first, cheaply.
3. **`GetKubernetesProvider` (188) / `GetKubernetesProviderAndPermissions` (264)** —
   highest volume but plausibly the most contention-sensitive of the
   provider-side queries; re-evaluate after Step 0.
4. **`CreateKubernetesResource` (137) / `ListKubernetesResourcesByTaskID` (562)** —
   lowest standalone suspicion; most likely to substantially improve from
   CN-5272 alone with no further work needed.
