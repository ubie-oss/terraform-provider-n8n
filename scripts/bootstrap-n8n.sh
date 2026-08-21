#!/usr/bin/env sh
# Bootstrap a local n8n instance for Terraform provider acceptance tests.
#
# Harness-only: uses the internal /rest API to create an owner and mint a
# Public API key. The Terraform provider must keep using /api/v1 + X-N8N-API-KEY.
#
# Usage:
#   eval "$(scripts/bootstrap-n8n.sh --export)"
#   scripts/bootstrap-n8n.sh            # prints N8N_ENDPOINT=... and N8N_API_KEY=...
#   scripts/bootstrap-n8n.sh --github-env  # appends to $GITHUB_ENV
#
# Env overrides:
#   N8N_BASE (default http://127.0.0.1:5678)
#   N8N_OWNER_EMAIL, N8N_OWNER_PASSWORD, N8N_OWNER_FIRST_NAME, N8N_OWNER_LAST_NAME

set -eu

N8N_BASE="${N8N_BASE:-http://127.0.0.1:5678}"
N8N_OWNER_EMAIL="${N8N_OWNER_EMAIL:-tf-acc@example.com}"
N8N_OWNER_PASSWORD="${N8N_OWNER_PASSWORD:-TfAccN8n1!}"
N8N_OWNER_FIRST_NAME="${N8N_OWNER_FIRST_NAME:-Tf}"
N8N_OWNER_LAST_NAME="${N8N_OWNER_LAST_NAME:-Acc}"
BROWSER_ID="${N8N_BROWSER_ID:-tf-acc-bootstrap-browser-id}"
WAIT_SECONDS="${N8N_BOOTSTRAP_WAIT_SECONDS:-120}"

MODE="print"
case "${1-}" in
--export) MODE="export" ;;
--github-env) MODE="github-env" ;;
-h | --help)
	sed -n '2,16p' "$0"
	exit 0
	;;
"") ;;
*)
	echo "unknown option: $1" >&2
	exit 2
	;;
esac

if ! command -v curl >/dev/null 2>&1; then
	echo "curl is required" >&2
	exit 1
fi
if ! command -v python3 >/dev/null 2>&1; then
	echo "python3 is required" >&2
	exit 1
fi

REPO_ROOT="$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)"
STATE_DIR="${REPO_ROOT}/.n8n-dev"
mkdir -p "${STATE_DIR}"
COOKIE_JAR="${STATE_DIR}/cookies.txt"
API_KEY_FILE="${STATE_DIR}/api-key"
rm -f "${COOKIE_JAR}"
# Pre-create the jar so curl writes into a 0600 file (session cookies).
(
	umask 077
	touch "${COOKIE_JAR}"
)
chmod 600 "${COOKIE_JAR}"

wait_ready() {
	i=0
	while [ "${i}" -lt "${WAIT_SECONDS}" ]; do
		if curl -sf "${N8N_BASE}/healthz/readiness" >/dev/null 2>&1; then
			return 0
		fi
		i=$((i + 1))
		sleep 1
	done
	echo "n8n did not become ready at ${N8N_BASE}/healthz/readiness within ${WAIT_SECONDS}s" >&2
	exit 1
}

# curl helpers: always send browser-id; cookie jar for session.
rest_post() {
	path="$1"
	body="$2"
	curl -sS -c "${COOKIE_JAR}" -b "${COOKIE_JAR}" \
		-H "Content-Type: application/json" \
		-H "browser-id: ${BROWSER_ID}" \
		-X POST "${N8N_BASE}/rest${path}" \
		-d "${body}" \
		-w "\n%{http_code}"
}

rest_get() {
	path="$1"
	curl -sS -c "${COOKIE_JAR}" -b "${COOKIE_JAR}" \
		-H "browser-id: ${BROWSER_ID}" \
		-X GET "${N8N_BASE}/rest${path}" \
		-w "\n%{http_code}"
}

# Split body and status from curl -w "\n%{http_code}" output.
split_response() {
	raw="$1"
	HTTP_BODY="$(printf '%s' "${raw}" | sed '$d')"
	HTTP_CODE="$(printf '%s' "${raw}" | tail -n 1)"
}

# Unwrap n8n RestController envelope {"data": ...} when present.
python_unwrap_data='
import json, sys
obj = json.load(sys.stdin)
if isinstance(obj, dict) and "data" in obj:
    json.dump(obj["data"], sys.stdout)
else:
    json.dump(obj, sys.stdout)
'

wait_ready

setup_body="$(OWNER_EMAIL="${N8N_OWNER_EMAIL}" OWNER_FIRST="${N8N_OWNER_FIRST_NAME}" OWNER_LAST="${N8N_OWNER_LAST_NAME}" OWNER_PASS="${N8N_OWNER_PASSWORD}" python3 -c '
import json, os
print(json.dumps({
  "email": os.environ["OWNER_EMAIL"],
  "firstName": os.environ["OWNER_FIRST"],
  "lastName": os.environ["OWNER_LAST"],
  "password": os.environ["OWNER_PASS"],
}))
')"

split_response "$(rest_post "/owner/setup" "${setup_body}")"
setup_code="${HTTP_CODE}"
setup_body_resp="${HTTP_BODY}"

if [ "${setup_code}" = "200" ]; then
	: # owner created; cookie jar populated
elif [ "${setup_code}" = "400" ]; then
	# Owner already exists — log in with the same credentials.
	login_body="$(OWNER_EMAIL="${N8N_OWNER_EMAIL}" OWNER_PASS="${N8N_OWNER_PASSWORD}" python3 -c '
import json, os
print(json.dumps({
  "emailOrLdapLoginId": os.environ["OWNER_EMAIL"],
  "password": os.environ["OWNER_PASS"],
}))
')"
	split_response "$(rest_post "/login" "${login_body}")"
	if [ "${HTTP_CODE}" != "200" ]; then
		echo "owner setup returned ${setup_code}; login failed with ${HTTP_CODE}: ${HTTP_BODY}" >&2
		exit 1
	fi
else
	echo "owner setup failed with HTTP ${setup_code}: ${setup_body_resp}" >&2
	exit 1
fi

split_response "$(rest_get "/api-keys/scopes")"
if [ "${HTTP_CODE}" != "200" ]; then
	echo "GET /rest/api-keys/scopes failed with HTTP ${HTTP_CODE}: ${HTTP_BODY}" >&2
	exit 1
fi
scopes_json="${HTTP_BODY}"

scopes_array="$(printf '%s' "${scopes_json}" | python3 -c "
${python_unwrap_data}
" | python3 -c '
import json, sys
scopes = json.load(sys.stdin)
if not isinstance(scopes, list) or not scopes:
    raise SystemExit("expected non-empty scopes list")
print(json.dumps(scopes))
')"

# Unique label so re-running bootstrap on the same instance does not collide.
KEY_LABEL="tf-acc-$(date +%s)-$$"
create_body="$(SCOPES_JSON="${scopes_array}" KEY_LABEL="${KEY_LABEL}" python3 -c '
import json, os
scopes = json.loads(os.environ["SCOPES_JSON"])
print(json.dumps({
  "label": os.environ["KEY_LABEL"],
  "scopes": scopes,
  "expiresAt": None,
}))
')"

split_response "$(rest_post "/api-keys" "${create_body}")"
if [ "${HTTP_CODE}" != "200" ]; then
	echo "POST /rest/api-keys failed with HTTP ${HTTP_CODE}: ${HTTP_BODY}" >&2
	echo "If this is CSRF/401, ensure browser-id is sent on setup/login and api-keys calls." >&2
	exit 1
fi

API_KEY="$(printf '%s' "${HTTP_BODY}" | python3 -c "
${python_unwrap_data}
" | python3 -c '
import json, sys
data = json.load(sys.stdin)
key = data.get("rawApiKey") or data.get("apiKey")
if not key or key.startswith("***") or (len(key) > 3 and key[0] == "*"):
    raise SystemExit("rawApiKey missing or redacted in create response")
print(key)
')"

printf '%s\n' "${API_KEY}" >"${API_KEY_FILE}"
chmod 600 "${API_KEY_FILE}"

N8N_ENDPOINT="${N8N_BASE}"
N8N_API_KEY="${API_KEY}"

case "${MODE}" in
export)
	# Quote for safe eval by the caller (make testacc-docker).
	printf "export N8N_ENDPOINT=%s\n" "$(printf '%s' "${N8N_ENDPOINT}" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read().rstrip("\n")))')"
	printf "export N8N_API_KEY=%s\n" "$(printf '%s' "${N8N_API_KEY}" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read().rstrip("\n")))')"
	;;
github-env)
	if [ -z "${GITHUB_ENV-}" ]; then
		echo "GITHUB_ENV is not set" >&2
		exit 1
	fi
	# Mask before writing so Actions logs do not leak the minted key.
	echo "::add-mask::${N8N_API_KEY}"
	{
		echo "N8N_ENDPOINT=${N8N_ENDPOINT}"
		echo "N8N_API_KEY=${N8N_API_KEY}"
	} >>"${GITHUB_ENV}"
	;;
print)
	printf 'N8N_ENDPOINT=%s\n' "${N8N_ENDPOINT}"
	printf 'N8N_API_KEY=%s\n' "${N8N_API_KEY}"
	;;
esac
