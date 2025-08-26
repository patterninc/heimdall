FROM golang:1.24.6

RUN apt-get update && apt-get install -y nodejs npm awscli jq

WORKDIR /go/src/github.com/patterninc/heimdall

COPY . .

# set config file
COPY configs/local.yaml /etc/heimdall/heimdall.yaml
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY configs/spark-application.yaml /etc/heimdall/spark-application.yaml

# build executables
RUN ./build.sh

CMD [ "/usr/local/bin/entrypoint.sh" ]
