#!/usr/bin/env bash
#
# Applies (or rolls back) idx_kubernetes_resources_kind_covering directly
# against the clouddriver database, per:
#   docs/CN-5272/kind-covering-index-rollout-runbook.md
#
# Intended usage: run from inside the throwaway mysql client pod, e.g.
#
#   kubectl run db-client --rm -it --image=mysql:8.0 --namespace=spinnaker -- bash
#
# then, inside that pod, either `kubectl cp` this file in from another
# terminal, or paste it directly via heredoc:
#
#   cat <<'SCRIPT_EOF' > /tmp/apply_kind_covering_index.sh
#   ...(paste this file's contents)...
#   SCRIPT_EOF
#   chmod +x /tmp/apply_kind_covering_index.sh
#
# Then run it:
#
#   DB_PASSWORD='...' /tmp/apply_kind_covering_index.sh
#
# The script is idempotent: if the index already exists it skips straight to
# verification instead of re-running the ALTER. Safe to re-run.
#
# Env vars (all but DB_PASSWORD have defaults matching the values discussed):
#   DB_HOST      (default: 10.190.1.30)
#   DB_USER      (default: clouddriver)
#   DB_NAME      (default: clouddriver)
#   DB_PASSWORD  (required - no default, never pass this as a CLI flag)
#   REPLICA_HOST (optional - if set, also verifies the index landed there)
#
# Usage:
#   DB_PASSWORD='...' ./apply_kind_covering_index.sh            # apply + verify
#   DB_PASSWORD='...' ./apply_kind_covering_index.sh --rollback # DROP INDEX
#   DB_PASSWORD='...' ./apply_kind_covering_index.sh --check-only  # no changes, just report status

set -euo pipefail

DB_HOST="${DB_HOST:-10.190.1.30}"
DB_USER="${DB_USER:-clouddriver}"
DB_NAME="${DB_NAME:-clouddriver}"
DB_PASSWORD="${DB_PASSWORD:-}"
REPLICA_HOST="${REPLICA_HOST:-}"

TABLE="kubernetes_resources"
INDEX_NAME="idx_kubernetes_resources_kind_covering"
INDEX_COLUMNS="kind,account_name,name,spinnaker_app"

MODE="apply"
for arg in "$@"; do
  case "$arg" in
    --rollback)   MODE="rollback" ;;
    --check-only) MODE="check-only" ;;
    *) echo "Unknown argument: $arg" >&2; exit 1 ;;
  esac
done

if [[ -z "$DB_PASSWORD" ]]; then
  echo "ERROR: DB_PASSWORD must be set (env var), e.g.:" >&2
  echo "  DB_PASSWORD='...' $0" >&2
  exit 1
fi

# Write credentials to a mysql defaults file instead of passing -p on the
# command line (which leaks the password into process listings/shell
# history). chmod 600 + auto-cleanup on exit.
CNF="$(mktemp)"
chmod 600 "$CNF"
trap 'rm -f "$CNF"' EXIT
cat > "$CNF" <<EOF
[client]
user=${DB_USER}
password=${DB_PASSWORD}
EOF

run_sql() {
  local host="$1" query="$2"
  mysql --defaults-extra-file="$CNF" -h "$host" -N -B "$DB_NAME" -e "$query"
}

echo "=== Pre-flight ==="

ROW_COUNT="$(run_sql "$DB_HOST" "SELECT COUNT(*) FROM ${TABLE};")"
echo "Row count: ${ROW_COUNT}"

ENGINE="$(run_sql "$DB_HOST" "SELECT ENGINE FROM information_schema.tables WHERE table_schema='${DB_NAME}' AND table_name='${TABLE}';")"
echo "Engine: ${ENGINE}"
if [[ "$ENGINE" != "InnoDB" ]]; then
  echo "ERROR: expected InnoDB, got '${ENGINE}'. ALGORITHM=INPLACE/LOCK=NONE assumptions may not hold. Aborting." >&2
  exit 1
fi

INDEX_EXISTS_COUNT="$(run_sql "$DB_HOST" "SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema='${DB_NAME}' AND table_name='${TABLE}' AND index_name='${INDEX_NAME}';")"
if [[ "$INDEX_EXISTS_COUNT" -gt 0 ]]; then
  INDEX_ALREADY_EXISTS=true
  echo "Index ${INDEX_NAME} already exists."
else
  INDEX_ALREADY_EXISTS=false
  echo "Index ${INDEX_NAME} does not exist yet."
fi

if [[ "$MODE" == "check-only" ]]; then
  echo "=== check-only mode: no changes made. Skipping to verification. ==="
elif [[ "$MODE" == "rollback" ]]; then
  if [[ "$INDEX_ALREADY_EXISTS" == "false" ]]; then
    echo "Index does not exist - nothing to roll back."
    exit 0
  fi
  echo "=== Rolling back: dropping ${INDEX_NAME} ==="
  time run_sql "$DB_HOST" "ALTER TABLE ${TABLE} DROP INDEX ${INDEX_NAME}, ALGORITHM=INPLACE, LOCK=NONE;"
  echo "Dropped."
  exit 0
else
  if [[ "$INDEX_ALREADY_EXISTS" == "true" ]]; then
    echo "=== Index already present - skipping ALTER, proceeding to verification. ==="
  else
    echo "=== Applying: creating ${INDEX_NAME} (${INDEX_COLUMNS}) ==="
    time run_sql "$DB_HOST" "ALTER TABLE ${TABLE} ADD INDEX ${INDEX_NAME} (${INDEX_COLUMNS}), ALGORITHM=INPLACE, LOCK=NONE;"
    echo "ALTER complete."

    echo "=== Refreshing planner statistics ==="
    run_sql "$DB_HOST" "ANALYZE TABLE ${TABLE};"
  fi
fi

echo
echo "=== Verification ==="

echo "--- SHOW INDEX (expect columns in order: kind, account_name, name, spinnaker_app) ---"
run_sql "$DB_HOST" "SELECT SEQ_IN_INDEX, COLUMN_NAME FROM information_schema.statistics WHERE table_schema='${DB_NAME}' AND table_name='${TABLE}' AND index_name='${INDEX_NAME}' ORDER BY SEQ_IN_INDEX;"

echo
echo "--- EXPLAIN: ListKubernetesClustersByApplication-equivalent ---"
EXPLAIN_1="$(run_sql "$DB_HOST" "EXPLAIN SELECT account_name, cluster FROM ${TABLE} WHERE spinnaker_app = 'verification-probe' AND kind IN ('Deployment','StatefulSet','ReplicaSet','Ingress','Service','DaemonSet') GROUP BY account_name, cluster;")"
echo "$EXPLAIN_1"
if echo "$EXPLAIN_1" | grep -q "$INDEX_NAME"; then
  echo "OK: planner is using ${INDEX_NAME}."
else
  echo "WARNING: planner did NOT report using ${INDEX_NAME} for this query. Investigate before relying on this index (see runbook Step 6)." >&2
fi

echo
echo "--- EXPLAIN: ListKubernetesClustersByFields-equivalent ---"
EXPLAIN_2="$(run_sql "$DB_HOST" "EXPLAIN SELECT account_name, kind, name, spinnaker_app FROM ${TABLE} WHERE kind IN ('Deployment','StatefulSet','ReplicaSet','Ingress','Service','DaemonSet') GROUP BY account_name, kind, name, spinnaker_app;")"
echo "$EXPLAIN_2"
if echo "$EXPLAIN_2" | grep -q "$INDEX_NAME"; then
  echo "OK: planner is using ${INDEX_NAME}."
else
  echo "WARNING: planner did NOT report using ${INDEX_NAME} for this query. Investigate before relying on this index (see runbook Step 6)." >&2
fi

if [[ -n "$REPLICA_HOST" ]]; then
  echo
  echo "=== Verifying replica (${REPLICA_HOST}) picked up the schema change (runbook Step 8) ==="
  REPLICA_INDEX_COUNT="$(run_sql "$REPLICA_HOST" "SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema='${DB_NAME}' AND table_name='${TABLE}' AND index_name='${INDEX_NAME}';" || echo "0")"
  if [[ "$REPLICA_INDEX_COUNT" -gt 0 ]]; then
    echo "OK: replica has the index."
  else
    echo "WARNING: replica does NOT yet have the index. If it's not actively replicating, it will need the ALTER applied manually (see runbook Step 8)." >&2
  fi
fi

echo
echo "=== Done ==="
