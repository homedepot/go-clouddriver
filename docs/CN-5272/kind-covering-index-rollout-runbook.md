# Runbook: rolling out `idx_kubernetes_resources_kind_covering` via manual DDL

JIRA: [CN-5272](https://thd.atlassian.net/browse/CN-5272)

This is the execution plan for Decision 3 in
[kind-casing-and-index-fix.md](kind-casing-and-index-fix.md) — creating
`idx_kubernetes_resources_kind_covering` directly against the production database
instead of letting `AutoMigrate` create it on pod startup, then removing it from
the Go struct so `AutoMigrate` never touches it again.

## Facts gathered during pre-flight (Step 0) that shape this plan

| Fact | Value | Why it matters |
|---|---|---|
| Row count | 630,613 | Small enough that `ALGORITHM=INPLACE, LOCK=NONE` is expected to complete in seconds to low single-digit minutes — see "Why we are not using gh-ost/pt-online-schema-change" below. |
| Engine | Cloud SQL for MySQL 8.0 | Confirms `ALGORITHM=INPLACE` secondary-index support and standard `information_schema`/`SHOW INDEX` verification syntax apply. |
| Read replica traffic | None live; one non-live replica exists | No live-replica lag to monitor during the DDL, but the non-live replica still needs its schema verified afterward (Step 8) so it doesn't silently drift and become a landmine if ever promoted. |
| Clouddriver app replicas | 8 | Confirms the `AutoMigrate`-race risk (8 pods potentially racing an `ALTER TABLE` on a rolling deploy) is real, not hypothetical — this is what Step 9 (removing the index from the Go struct) actually eliminates. |

### Why we are not using gh-ost / pt-online-schema-change

Those tools exist to make DDL non-blocking on **huge** tables (tens of millions+
rows), where even `LOCK=NONE` can leave a long-running operation with cumulative
risk, or where the brief metadata lock at the very start/end of the DDL is still
unacceptable. Adding a secondary (non-unique) index via
`ALGORITHM=INPLACE, LOCK=NONE` in InnoDB does **not** rebuild the table — MySQL
only builds the new index structure and merges it in, with full read/write
concurrency the entire time except a sub-second metadata lock at the very
start/end. At 630,613 rows, this is expected to be fast. Standing up gh-ost/pt-osc
(binlog `ROW` format, elevated privileges, tooling not currently present in this
infra per the availability check) would be real setup work for a problem this
table's size doesn't have. If a future migration needs to run against a
significantly larger table, revisit gh-ost/pt-osc at that point — this decision is
specific to this table's current size.

---

## Scripted execution (Steps 4-8, and rollback)

[`scripts/apply_kind_covering_index.sh`](scripts/apply_kind_covering_index.sh)
automates the DDL application (Step 4) and verification (Steps 5-8) into one
idempotent script — safe to re-run, since it checks whether the index already
exists before doing anything. It does **not** replace Steps 1-3 (generating/
confirming the DDL, rehearsing in staging, scheduling the window/backup) or
Step 9-10 (removing the index from the Go struct, documenting the
new-environment gap) — those still require human judgment/timing decisions.

**Using it, from the throwaway client pod you already use to reach this DB:**

```bash
kubectl run db-client --rm -it --image=mysql:8.0 --namespace=spinnaker -- bash
```

Then, inside that pod's shell, paste the script directly via heredoc (no
`kubectl cp` / second terminal needed):

```bash
cat <<'SCRIPT_EOF' > /tmp/apply_kind_covering_index.sh
# (paste the contents of docs/CN-5272/scripts/apply_kind_covering_index.sh here)
SCRIPT_EOF
chmod +x /tmp/apply_kind_covering_index.sh
```

Then run it (password via env var, never as a `-p` flag, so it never lands in
shell history or `ps`):

```bash
# Dry run first - reports current state, makes no changes:
DB_PASSWORD='...' /tmp/apply_kind_covering_index.sh --check-only

# Apply for real (Step 4) and auto-run verification (Steps 5-7):
DB_PASSWORD='...' /tmp/apply_kind_covering_index.sh

# Also verify the non-live replica (Step 8), if you can reach it from this pod:
DB_PASSWORD='...' REPLICA_HOST='<replica-ip>' /tmp/apply_kind_covering_index.sh

# Roll back (drops the index):
DB_PASSWORD='...' /tmp/apply_kind_covering_index.sh --rollback
```

Defaults match what's been discussed: `DB_HOST=10.190.1.30`, `DB_USER=clouddriver`,
`DB_NAME=clouddriver` — override any of them as env vars if they're wrong.

What it does, mapped to the manual steps below (read those for the *why*; the
script is the *how*):
- Pre-checks: row count, confirms `ENGINE=InnoDB`, checks whether the index
  already exists (skips the `ALTER` if so, instead of erroring).
- Applies the DDL from Step 4, timed.
- Runs `ANALYZE TABLE` (Step 7).
- Verifies index structure via `information_schema.statistics` (Step 5).
- Runs `EXPLAIN` against both target queries and greps for the index name to
  confirm the planner is actually using it (Step 6) — prints a `WARNING` rather
  than failing silently if it isn't.
- Optionally checks the replica (Step 8) if `REPLICA_HOST` is set.

The manual SQL in Steps 4-8 below is kept as the documented reference for
what the script runs and why — useful if you need to run a step by hand,
diagnose a `WARNING` from the script, or the script itself needs updating.

---

## Step 1 — Generate the exact DDL from the Go struct definition

Don't hand-write the index DDL from memory. Generate it once against a scratch
database so the column list and order exactly match what
`internal/kubernetes/resource.go`'s `gorm` tag currently specifies, before that tag
is removed in Step 9.

```go
// scratch program, run once against a throwaway SQLite or local MySQL DB
package main

import (
	"fmt"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // prints generated SQL
	})
	db.AutoMigrate(&kubernetes.Resource{})
}
```

Confirm the printed `CREATE INDEX`/index clause matches:

```sql
ALTER TABLE kubernetes_resources
  ADD INDEX idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app),
  ALGORITHM=INPLACE, LOCK=NONE;
```

This is the exact statement used in every step below. If the scratch run produces
a different column order, use that instead — the code, not this doc, is the source
of truth for what GORM would have created.

---

## Step 2 — Rehearse in a non-production environment

Run the Step 1 statement against staging (or any non-prod MySQL instance with a
similar `kubernetes_resources` shape) first.

```bash
cloud-sql-proxy <staging-instance-connection-name> &
mysql -h 127.0.0.1 -u <user> -p <database> <<'EOF'
ALTER TABLE kubernetes_resources
  ADD INDEX idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app),
  ALGORITHM=INPLACE, LOCK=NONE;
EOF
```

Record:
- Wall-clock duration (`time` the command, or check MySQL's own reported execution time).
- Any CPU/IOPS spike on the staging instance's Cloud SQL monitoring dashboard during the run.

Given the row count, expect this to be fast (seconds to low minutes) with no
visible impact — but confirm it rather than assume it, since this is the number
you'll compare against when running it for real.

---

## Step 3 — Schedule the production change

1. Pick a low-traffic window. `LOCK=NONE` should mean no DML blocking, but treat
   the staging timing as an estimate, not a guarantee — don't run this during a
   peak deploy window.
2. Take an on-demand Cloud SQL backup immediately before the change:
   ```bash
   gcloud sql backups create --instance=<prod-instance-name>
   ```
   This gives you a rollback point independent of the DDL succeeding or not.
3. Notify whoever is on call for Spinnaker/clouddriver, including the expected
   duration from Step 2.

---

## Step 4 — Run the DDL against production

```bash
cloud-sql-proxy <prod-instance-connection-name> &
mysql -h 127.0.0.1 -u <user> -p <database> <<'EOF'
ALTER TABLE kubernetes_resources
  ADD INDEX idx_kubernetes_resources_kind_covering (kind, account_name, name, spinnaker_app),
  ALGORITHM=INPLACE, LOCK=NONE;
EOF
```

Watch Cloud SQL's monitoring (CPU, IOPS) live while it runs. If something looks
wrong mid-flight, `SHOW PROCESSLIST` in another session to find the ALTER's
connection ID and `KILL <id>` — `LOCK=NONE` operations are designed to be safely
abortable without corrupting the table, but confirm this holds for MySQL 8.0
before relying on it under pressure.

---

## Step 5 — Verify the index exists and is structured correctly

```sql
SHOW INDEX FROM kubernetes_resources WHERE Key_name = 'idx_kubernetes_resources_kind_covering';
```

Confirm the column order in the output matches `(kind, account_name, name, spinnaker_app)`
exactly — order matters for whether the index can serve as a covering index for
the target queries.

---

## Step 6 — Verify the query planner actually uses it

```sql
EXPLAIN SELECT account_name, cluster FROM kubernetes_resources
  WHERE spinnaker_app = 'some-real-app-name'
    AND kind IN ('Deployment','StatefulSet','ReplicaSet','Ingress','Service','DaemonSet')
  GROUP BY account_name, cluster;
```

Look for:
- `key: idx_kubernetes_resources_kind_covering`
- `Extra: Using index` (confirms an index-only scan, not a scan that still hits the base table)

If `EXPLAIN` instead shows a different key or `Using where` without `Using index`,
stop and investigate before moving on — the index exists but isn't being used as
intended, and Step 9 (removing it from the Go struct) should not proceed until
this is resolved.

---

## Step 7 — Refresh planner statistics

```sql
ANALYZE TABLE kubernetes_resources;
```

Ensures the optimizer's row-count/cardinality estimates for the new index are
current immediately after creation, rather than relying on stale statistics for
early queries against it.

---

## Step 8 — Verify the non-live replica picks up the schema change

Since Cloud SQL replicates DDL through the binlog like any other statement, the
non-live replica should receive this `ALTER TABLE` automatically. Confirm it did,
so it doesn't silently drift from the primary and become a landmine if it's ever
promoted (failover, DR test, etc.):

```sql
-- run against the replica
SHOW INDEX FROM kubernetes_resources WHERE Key_name = 'idx_kubernetes_resources_kind_covering';
```

If the replica is not currently replicating (fully stopped, not just "not serving
traffic"), this index will need to be applied to it separately using the same
statement from Step 1 — confirm its replication status before assuming it will
catch up on its own.

---

## Step 9 — Remove the index from the Go struct

Only after Steps 5–8 all confirm the index is live, correctly structured, used by
the planner, and present on the replica:

1. In `internal/kubernetes/resource.go`, remove
   `index:idx_kubernetes_resources_kind_covering,priority:N` from the `AccountName`,
   `Name`, `Kind`, and `SpinnakerApp` field tags, leaving `account_name_kind_name_spinnaker_app_idx`
   and `kind_idx` (the pre-existing indexes) untouched.
2. Update the doc comment above the `Resource` struct to state this index is
   managed manually per this runbook, not by `AutoMigrate`, and link to this file.
3. This is the change that actually eliminates the 8-replica `AutoMigrate` race —
   once the tag is gone, no pod's startup will attempt to create, verify, or touch
   this index at all.

---

## Step 10 — Close the "new environment" gap

Because `AutoMigrate` no longer creates this index, any brand-new environment
(new staging cluster, DR restore into an empty schema, a new tenant DB) will not
get it automatically. Add the Step 1 statement to whatever provisioning process
stands up a fresh `kubernetes_resources` table — at minimum, a note in this
runbook that says "run Step 1's statement by hand after first deploy to any new
environment" — so this doesn't become a silent, undocumented manual step that
only lives in one person's memory.

---

## Step 11 — Monitor after rollout

- Track latency on `ListKubernetesClustersByApplication` /
  `ListKubernetesClustersByFields` (whatever dashboard surfaces clouddriver query
  timing) for the expected improvement.
- Watch write-side metrics (`INSERT` latency, IOPS) over the following few days to
  confirm the index's write-amplification cost is acceptable at this table's
  actual write rate — this is also the evidence base for whether
  `idx_kubernetes_resources_app_covering` (deferred in Decision 2) ever becomes
  necessary.

---

## Rollback plan

If the index needs to be removed after the fact (unexpected write-latency
regression, planner picking a worse plan for some other query, etc.):

```sql
ALTER TABLE kubernetes_resources
  DROP INDEX idx_kubernetes_resources_kind_covering,
  ALGORITHM=INPLACE, LOCK=NONE;
```

Dropping a secondary index is always fast and low-risk (no table rebuild). If this
happens, also revert Step 9's code change so the doc comment and struct tags stay
truthful about what's actually in the database.
