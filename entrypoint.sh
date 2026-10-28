#!/bin/bash

# start heimdall UI
(cd ./web && ./node_modules/.bin/next start &)

# start heimdall
(cd ./dist && ./deploydb ../assets/databases/heimdall/build/heimdall.lst && ./heimdall)
