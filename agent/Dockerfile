FROM golang:1.14 AS build
WORKDIR /go/src/github.com/blaqkube/mysql-operator/agent
COPY * ./
RUN go test -cover ./...

ENV CGO_ENABLED=0
RUN go build -a -installsuffix cgo -o mysql-agent .

FROM mysql:8.0.20
RUN apt update && \
    apt install -y ca-certificates
COPY --from=build /go/src/github.com/blaqkube/mysql-operator/agent/mysql-agent ./
EXPOSE 8080/tcp
ENTRYPOINT ["./mysql-agent"]
