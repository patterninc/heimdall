#!/bin/bash

# start heimdall UI
(cd ./web && pnpm start &)

# start heimdall
(cd ./dist && ./deploydb ../assets/databases/heimdall/build/heimdall.lst && ./heimdall)
