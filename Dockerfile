FROM golang:1.24.3

RUN apt-get update && apt-get install -y nodejs npm awscli jq

WORKDIR /go/src/github.com/patterninc/heimdall

COPY . .

# set config file
COPY configs/local.yaml /etc/heimdall/heimdall.yaml
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# build executables
RUN go version && go mod tidy
RUN ./build.sh

CMD [ "/usr/local/bin/entrypoint.sh" ]
