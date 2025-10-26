#! /bin/bash

source ./dev_env/load_env_vars.sh && /$(which dlv) exec \
--listen=:2345 \
--headless=true \
--api-version=2 \
--accept-multiclient \
--log \
--continue \
./tmp/dyutas-auth