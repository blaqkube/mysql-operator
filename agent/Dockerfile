FROM golang:1.14 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .

ENV CGO_ENABLED=0
RUN go build -a -installsuffix cgo -o mysql-agent .

FROM scratch AS runtime
COPY --from=build /go/src/mysql-agent ./
EXPOSE 8080/tcp
ENTRYPOINT ["./mysql-agent"]
