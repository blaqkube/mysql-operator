# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY agent/ agent/

# Build
RUN mkdir /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o /app/manager main.go

FROM alpine:3.13

ARG agent_version=0
ENV AGENT_VERSION=$agent_version

WORKDIR /app
COPY --from=builder /app/manager .
RUN chown -R 65532 /app

ENTRYPOINT ["/app/manager"]
