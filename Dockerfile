# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o prometheus-test-data ./cmd/main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
LABEL vendor="IBM" \
  name="prometheus-test-data" \
  org.label-schema.vendor="IBM" \
  org.label-schema.name="ibm fusion" \  
  org.label-schema.vcs-url="https://github.com/shdn-ibm/prometheus-test-data" \
  org.label-schema.schema-version="1.0.0"

WORKDIR /

COPY --from=builder /workspace/prometheus-test-data .

ENTRYPOINT ["/prometheus-test-data"]