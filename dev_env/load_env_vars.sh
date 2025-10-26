#! /bin/bash

if [ -f .env ]; then
    set -o allexport
    source ./.env
    set +o allexport
    # Export all variables explicitly
    export $(grep -v '^#' .env | xargs)
else
    echo "Error: auth service .env file not found"
    exit 1
fi

DYUTAS_ENV_PATH="$DYUTAS_ENV_PATH"
if [ -z "$DYUTAS_ENV_PATH" ]; then
    echo "Error: DYUTAS_ENV_PATH env var is not set"
    exit 1
fi

if [ -f "$DYUTAS_ENV_PATH" ]; then
    set -o allexport
    source "$DYUTAS_ENV_PATH"
    set +o allexport
    # Export all variables explicitly
    export $(grep -v '^#' "$DYUTAS_ENV_PATH" | xargs)
else
    echo "Error: DYUTAS_ENV file not found(env var DYUTAS_ENV_PATH should be set) from: $DYUTAS_ENV_PATH"
    exit 1
fi
