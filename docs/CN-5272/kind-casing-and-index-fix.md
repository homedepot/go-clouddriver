# Decision Record: `kubernetes_resources.kind` query fix, index rollout, and casing normalization

Tracked internally as CN-5272.

## Context

`go-clouddriver` is THD's custom replacement for Spinnaker's clouddriver — it is the
component that talks to Kubernetes clusters and is the source of truth clouddriver's
callers (Deck's Clusters tab, Orca's target-resolution logic for
disable/rollback/red-black stages) rely on for "what clusters/resources exist for
this Spinnaker application."

Two queries backing that data — `ListKubernetesClustersByApplication` and
`ListKubernetesClustersByFields` (`internal/sql/client.go`) — were doing:

```sql
WHERE UPPER(kind) IN ('DEPLOYMENT', 'STATEFULSET', 'REPLICASET', 'INGRESS', 'SERVICE', 'DAEMONSET')
```

against `kubernetes_resources`, a table with 300k+ rows, at 1-5s per call. This
document records the decisions made while fixing that, and *why* — so the next
person touching this code understands the tradeoffs instead of just the diff.

Because this table underpins live cluster visibility and (indirectly) pipeline
target resolution, the standard held throughout was: **don't ship a fix whose
correctness depends on an unstated property of the current environment** (e.g. "it
works because MySQL's collation happens to be case-insensitive today"). If a
property must hold for correctness, it must be enforced by the code, not assumed.

---

## Decision 1 — Remove `UPPER(kind)`, compare `kind` directly

**What:** `WHERE UPPER(kind) IN (...)` → `WHERE kind IN (?)` against a `clusterKinds`
slice of canonical PascalCase literals (`"Deployment"`, `"StatefulSet"`, etc.).

**Why:** Wrapping a column in a function (`UPPER(kind)`) makes the predicate
non-sargable — MySQL can't use a B-tree index on `kind` and falls back to a full
table/index scan, which is exactly the 1-5s cost being measured. Comparing `kind`
directly restores sargability and lets an index on `kind` be used as a range/ref
scan instead.

**Risk this introduces:** correctness now depends on every row's `kind` value
actually being (or collating as) the canonical PascalCase form used in the literal
list. See Decision 4 — this is the one open question this change created, and it's
closed by normalizing at write time rather than leaning on collation.

---

## Decision 2 — Add one covering index now (`idx_kubernetes_resources_kind_covering`), defer the second

**What:** `internal/kubernetes/resource.go` adds
`idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app)`,
a covering index for the two hot queries above. A second index,
`idx_kubernetes_resources_app_covering (spinnaker_app, account_name, cluster, kind)`,
was considered and **explicitly not added**.

**Why:**
- The `kind_covering` index is the direct fix for the measured bottleneck query — it
  lets both queries run as index-only scans instead of scanning the table.
- `kubernetes_resources` is written to continuously by a multi-tenant clouddriver
  polling many clusters — every additional index taxes every `INSERT`/`UPDATE` with
  write amplification (extra CPU/IOPS on Cloud SQL). Adding a second wide composite
  index speculatively, without slow-query-log evidence it's needed, trades a
  guaranteed ongoing write cost for a hypothetical read benefit.
- Standard operational practice for a write-heavy production table: ship the index
  that's proven necessary, measure, and only add more if the evidence says so.

**How to apply going forward:** if slow-query logs later show
`ListKubernetesAccountsBySpinnakerApp`/application-scoped lookups are still slow,
revisit adding `idx_kubernetes_resources_app_covering` at that point — not before.

---

## Decision 3 — How the index actually gets created in production (resolved — see internal runbook)

**What:** Both indexes (in Decision 2, and generally) are declared via GORM struct
tags and created by `db.AutoMigrate(...)` in `sql.Client.Connect()`, which runs on
**every pod's startup**.

**Why this needs a real decision before shipping:** `AutoMigrate` running a blind
`ALTER TABLE ADD INDEX` against a live, 300k+-row, write-heavy table on every
replica's boot is an operational risk, not just a code concern:
- Multiple clouddriver replicas could race through `AutoMigrate` concurrently on a
  rolling deploy; GORM's `HasIndex` + `CreateIndex` isn't atomic across processes,
  so this could surface as a "duplicate key name" error on one replica.
- An uncontrolled `ALTER TABLE` on a large table can cause lock contention/caching
  lag depending on the DB engine and current load, even with mostly-online DDL in
  modern MySQL.

**Resolution:** the index tag **stays declared in the Go struct** — `AutoMigrate`
remains the source of truth and continues to be how this index gets created by
default. This matters for anyone else standing up a fresh install of this
project (a new environment, a new tenant DB, a first-time deploy): they get this
index automatically, with no manual step required, the same as any other index
in this codebase.

The manual, directly-run `ALTER TABLE ... ALGORITHM=INPLACE, LOCK=NONE` DDL is a
**pre-emptive, conditional step** for a specific operational situation, not a
replacement for `AutoMigrate`: run it by hand, once, *before* deploying this
change, **only if both of the following are true**:
- More than one replica of this service is running against the same database
  (the race described above requires concurrent `AutoMigrate` calls to actually
  matter — a single-replica deploy has nothing to race against).
- The table has grown large enough that a live `ALTER TABLE` against it is a
  real operational concern for your environment (no fixed row count applies
  universally here — it depends on your DB's resources and tolerance for a
  brief DDL operation; teams running a much larger `kubernetes_resources` table
  than the one this decision was made against should evaluate their own
  threshold rather than assume this one's numbers apply).

When both conditions hold, applying the index manually first means every
replica's subsequent `AutoMigrate` call finds the index already present via its
`HasIndex` check and skips `CreateIndex` — the race is avoided without touching
the code, and the index stays declared in the struct for the benefit of every
other environment. gh-ost/pt-online-schema-change were evaluated and **not
used** for this: at this table's actual size (630,613 rows as of this writing),
a secondary-index add via `ALGORITHM=INPLACE` doesn't rebuild the table and is
expected to complete in seconds to low minutes, which is well within what a
plain `ALTER TABLE` handles safely — those tools solve a problem (long-running
blocking DDL on huge tables) this table doesn't have at its current scale. The
full step-by-step execution plan (including production connection details) is
maintained internally, not in this public repo.

---

## Decision 4 — Normalize `kind` at write time instead of relying on DB collation

**What:** Added `kubernetes.Client.GVKForKind(kind string) (schema.GroupVersionKind, error)`
(`internal/kubernetes/client.go`), a thin wrapper around the REST mapper's
`KindFor` resolution — the same resolution `Get()` and
`DeleteResourceByKindAndNameAndNamespace()` already perform internally. Then:

| File | Before | After | Why this approach |
|---|---|---|---|
| `internal/api/core/kubernetes/delete.go` (both dynamic/static and label-selector branches) | `Kind: kind` (raw, parsed from `ManifestName`/`dm.Kinds` — often lowercase) | `Kind: gvk.Kind` via new `GVKForKind` call | List-response items from `ListByGVR` do **not** reliably carry `Kind`/`apiVersion` — verified in `client-go`'s `dynamic/simple.go`: `List()` decodes the raw API response with no backfilling, and the Kubernetes API convention is that embedded list items omit `TypeMeta`. So `item.GetKind()` was not a safe source here; the REST-mapper resolution already computed elsewhere in this code path was the only reliable canonical source. |
| `internal/api/core/kubernetes/enable.go`, `disable.go` | `Kind: kind` (same raw source) | `Kind: target.GetKind()` | `target` is a single-object `Get()` result, which — unlike list items — reliably carries `TypeMeta`. No new client method needed; the canonical value was already sitting unused in a variable already in scope. |
| `internal/api/core/kubernetes/runjob.go` | `Kind: "job"` (hardcoded, wrong casing) | `Kind: meta.Kind` | `meta.Kind` is already set from the resolved GVK inside `Apply`/`Replace` (`internal/kubernetes/client.go`) — the correct value was already computed and discarded. |

**Why normalize at write time instead of leaving `kind IN (?)` to rely on MySQL's
current case-insensitive collation:**
- Collation is a property of the *database*, not the *code*. Today's MySQL default
  (`utf8mb3_unicode_ci`, confirmed against the live DB) is case-insensitive, so the
  query works today — but that's the DB masking a data-quality gap, not the code
  guaranteeing correctness.
- `internal/kubernetes/cluster.go`'s `Cluster()` function already self-normalizes
  the derived cluster-name string; the raw `kind` column itself had no equivalent
  protection before this change.
- Collation-dependent correctness is **already reachable today**, not
  hypothetical: local/dev runs on SQLite (case-sensitive by default — it's the
  fallback dialect when `DB_HOST`/`DB_NAME`/etc. are unset), so a developer testing
  `delete`/`enable`/`disable` locally could already reproduce clusters silently
  vanishing from `ListKubernetesClustersByApplication`/`ByFields` before this fix.
- If MySQL collation is ever hardened (a common compliance ask) or this DB is ever
  migrated to Postgres, a fix that relies on collation breaks silently in
  production with no error — just quietly incomplete cluster listings, which is the
  worst failure mode for a clouddriver replacement (wrong-but-quiet, not
  loud-and-caught).
- Concretely, the risk was **not** hypothetical for delete/enable/disable
  specifically: those paths derive `kind` from `ManifestName`, whose casing follows
  whatever Orca/Spinnaker convention sends it (existing tests before this fix used
  `"clusterRole"`, `"replicaset"`, `"deployment"` — all non-canonical), and those
  paths insert new, real rows into this append-only history table. A cluster whose
  only surviving rows came from a `disable`/`enable`/`delete` operation (a routine
  day-2 operation, not an edge case) could disappear entirely from the Clusters tab
  on any case-sensitive backend.

**Test coverage added to prove this, not just assert it:** each of
`delete_test.go`, `enable_test.go`, `disable_test.go`, `runjob_test.go` now has a
test that deliberately sets a lowercase/mismatched input casing and asserts the
*persisted* `Resource.Kind` is canonical PascalCase — proving the fix is
casing-input-independent, not just checking the happy path where input and output
already agree.

---

## Decision 5 — Removed a test that couldn't prove what it claimed

**What:** Deleted a test in `internal/sql/client_test.go` titled "kind casing in the
data differs from the canonical list."

**Why:** The test used `sqlmock`, which doesn't execute real SQL — it matches the
query text/args and then returns whatever rows you told it to return,
unconditionally. The test could not actually prove `kind IN ('Deployment', ...)`
matches a row stored as `'deployment'` (its own comment admitted as much). All it
verified was that Go-side row scanning doesn't care about the casing of *already
returned* data — never in doubt, and already covered by the adjacent "it succeeds"
test. Keeping a test whose name implied "we've verified the collation-dependent
behavior" when it verified no such thing would have been actively misleading to a
future reader. Per direction at the time: only keep tests that address what's
actually implemented now; add real regression coverage later if/when a test can
genuinely exercise the risk (e.g. an integration test against a real case-sensitive
backend).

---

## Summary of what's now true

- The query is sargable and index-backed (Decision 1 + 2).
- `kubernetes_resources.kind` is normalized to canonical PascalCase at every write
  path, by construction, not by accident of collation (Decision 4).
- The `kind IN (?)` comparison is therefore correct on **any** database/collation —
  MySQL with any collation, Postgres, or SQLite — not just today's specific
  production configuration.
- The index rollout mechanism (Decision 3) is resolved: the index stays declared
  in code via `AutoMigrate` for every environment by default, with a
  conditional manual pre-creation step (tracked in an internal runbook) only
  for deployments running multiple replicas against a database large enough
  that the risk applies. Still open, tracked separately: the longer-term
  `utf8mb3` → `utf8mb4` charset migration (unrelated technical debt, not a
  blocker for this work).
