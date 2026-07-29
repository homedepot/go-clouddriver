# Decision Record: `kubernetes_resources.kind` query fix, index evaluation, and casing normalization

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

## Decision 2 — Add a covering index? Tried, tested against real data, rolled back

**What was originally proposed:** `internal/kubernetes/resource.go` added
`idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app)`,
a covering index intended for the two hot queries above. A second index,
`idx_kubernetes_resources_app_covering (spinnaker_app, account_name, cluster, kind)`,
was considered and not added, on the reasoning that it should only be added if
evidence showed it was needed.

**What we did to validate it:** before shipping this permanently, we tested it
against a real, production-shaped test database (630k+ rows) rather than
relying on `EXPLAIN`'s estimated cost alone:
- `ListKubernetesClustersByFields`: real wall-clock timing was **identical**
  with the index present (1.39s) and absent (1.37s), same 90,103-row result.
  Forcing the planner to use the index (`FORCE INDEX`) made the plan *worse* —
  same row count scanned, plus an added `Using temporary` step the planner
  otherwise avoided by choosing the pre-existing
  `account_name_kind_name_spinnaker_app_idx` instead (which already contains
  all four selected/grouped columns and served as a covering index on its own,
  once `UPPER(kind)` was removed).
- `ListKubernetesClustersByApplication`: real wall-clock timing improved 5x
  (0.30s → 0.06s) after the query rewrite — but the query plan showed this
  came entirely from an existing index on `(spinnaker_app, kind)` (tracked
  separately, out of scope for this change — see below), not from
  `idx_kubernetes_resources_kind_covering`, which the planner never chose for
  this query either.

**Resolution: the index was removed.** Across both queries it was built for,
real data showed no measurable benefit — the entire performance improvement
came from Decision 1 (making `kind` sargable) combined with indexes that
already existed. Shipping an index that provides no demonstrated benefit would
mean paying its write-amplification cost (extra CPU/IOPS on every
`INSERT`/`UPDATE` to this write-heavy table, in every environment,
indefinitely) for nothing — exactly the tradeoff this decision originally said
to avoid "without slow-query-log evidence it's needed." The evidence came in,
and it said no.

**How to apply going forward:** if a future slow-query-log shows a genuine gap
not covered by existing indexes, design and validate a new index against real
data the same way — don't ship on `EXPLAIN`'s estimate alone; confirm the
planner actually chooses it and that real timing improves.

---

## Decision 3 — Retired (was: how the index gets created in production)

This decision existed to solve a rollout problem for the index proposed in
Decision 2 (avoiding an `AutoMigrate` race across replicas on a large table).
Since that index was removed entirely rather than shipped, **there is nothing
to roll out, and this decision no longer applies.** No index-creation DDL
needs to run anywhere, manually or via `AutoMigrate`, as a result of this
change.

The manual DDL runbook and its script that were drafted for this rollout
(maintained internally, not in this public repo) are now obsolete and should
be retired/removed rather than followed — they document a rollout for an
index this change no longer ships.

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

- The query is sargable, and real testing confirms pre-existing indexes
  (`account_name_kind_name_spinnaker_app_idx` and a separately-tracked index on
  `(spinnaker_app, kind)`) are sufficient to serve it well — no new index ships
  as part of this change (Decision 1 + 2).
- `kubernetes_resources.kind` is normalized to canonical PascalCase at every write
  path, by construction, not by accident of collation (Decision 4).
- The `kind IN (?)` comparison is therefore correct on **any** database/collation —
  MySQL with any collation, Postgres, or SQLite — not just today's specific
  production configuration.
- No index rollout is needed (Decision 3 is retired along with the index it
  existed to roll out). Still open, tracked separately: the longer-term
  `utf8mb3` → `utf8mb4` charset migration (unrelated technical debt, not a
  blocker for this work), and formalizing `spinnaker_app_kind_idx` into
  `AutoMigrate` (separate story, see below).

## Tracked separately, out of scope for this fix

- **`spinnaker_app_kind_idx (spinnaker_app, kind)` on `kubernetes_resources`** —
  discovered during validation testing to be present consistently across all
  environments (LLC and prod), and to be load-bearing for
  `ListKubernetesClustersByApplication`'s real performance (it's what the query
  planner actually uses; see validation notes). It is not declared in
  `resource.go`, not managed by `AutoMigrate`, and not documented in the
  README's "MySQL Indexes and Cleanup" section (verified - only five indexes
  are documented there, and this isn't one of them). Deliberately not touched
  as part of this work; a separate story will formalize it (add to
  `AutoMigrate`/document it) rather than folding it into this fix's scope.
