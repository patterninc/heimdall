#!/bin/bash

# start heimdall UI
(cd ./web && HOSTNAME=127.0.0.1 PORT=3000 node server.js &)

# start heimdall
(cd ./dist && ./deploydb ../assets/databases/heimdall/build/heimdall.lst && ./heimdall)
