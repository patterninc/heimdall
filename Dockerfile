FROM golang:1.24.6

RUN apt-get update && apt-get install -y nodejs npm awscli jq
RUN apt-get update && apt-get install -y --no-install-recommends build-essential binutils binutils-gold pkg-config && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/patterninc/heimdall

COPY . .

# set config file
COPY configs/local.yaml /etc/heimdall/heimdall.yaml
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
# build executables
RUN ./build.sh

CMD [ "/usr/local/bin/entrypoint.sh" ]