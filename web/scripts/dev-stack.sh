#!/bin/sh
# Starts postgres + the heimdall API from the repo's docker-compose.yaml, waits for the API,
# then runs `next dev` on :4000 against it. `--down` stops and removes the containers.
set -eu

WEB_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$WEB_DIR/../docker-compose.yaml"
API_URL=${API_URL:-http://localhost:9090}
TIMEOUT=${DEV_STACK_TIMEOUT:-180}

BUILD_FLAG=
DOWN=false
VOLUMES_FLAG=
for arg in "$@"; do
  case "$arg" in
    --build) BUILD_FLAG=--build ;;
    --down) DOWN=true ;;
    --volumes) VOLUMES_FLAG=--volumes ;;
    *) echo "unknown flag: $arg (use --build, --down, --volumes)" >&2; exit 1 ;;
  esac
done

if ! docker info > /dev/null 2>&1; then
  echo "dev:stack: Docker is not running or not installed. Start Docker and try again." >&2
  exit 1
fi

# docker-compose.yaml passes these through to the API container; default them to empty so
# compose does not warn when they are unset (only the AWS-backed plugins need them).
export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-}"
export AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN:-}"
export AWS_REGION="${AWS_REGION:-}"

if [ "$DOWN" = true ]; then
  docker compose -f "$COMPOSE_FILE" down $VOLUMES_FLAG
  exit 0
fi

echo "dev:stack: starting postgres and heimdall with docker compose..."
docker compose -f "$COMPOSE_FILE" up -d $BUILD_FLAG postgres heimdall

echo "dev:stack: waiting up to ${TIMEOUT}s for $API_URL/api/v1/commands..."
if ! node -e '
  const url = process.argv[1] + "/api/v1/commands"
  const deadline = Date.now() + Number(process.argv[2]) * 1000
  const probe = async () => {
    try {
      const res = await fetch(url, { signal: AbortSignal.timeout(2000) })
      if (res.ok) return process.exit(0)
    } catch {}
    if (Date.now() > deadline) return process.exit(1)
    setTimeout(probe, 1000)
  }
  probe()
' "$API_URL" "$TIMEOUT"; then
  echo "dev:stack: the API did not answer in ${TIMEOUT}s. Check the logs:" >&2
  echo "  docker compose -f $COMPOSE_FILE logs heimdall" >&2
  exit 1
fi

echo "dev:stack: API is up. Starting next dev on http://localhost:4000 (API_URL=$API_URL)"
cd "$WEB_DIR"
export API_URL
PORT=4000 exec ./node_modules/.bin/next dev
