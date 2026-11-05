#!/bin/bash

# start heimdall UI
(cd ./web && ./node_modules/.bin/next start &)

# start heimdall; exec replaces this shell so heimdall is PID 1 and receives the container's SIGTERM
cd ./dist && ./deploydb ../assets/databases/heimdall/build/heimdall.lst && exec ./heimdall
