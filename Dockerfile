FROM golang:1.17 as builder

LABEL maintainer="mail@host.tld"
LABEL stage=builder

# Set up execution environment in container's GOPATH
WORKDIR /go/src/app/cmd

# Copy relevant folders into container
COPY ./cmd /go/src/app/cmd
COPY ./internal /go/src/app/internal
COPY ./pkg /go/src/app/pkg
COPY ./tools /go/src/app/tools
COPY ./go.mod /go/src/app/go.mod

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

### NOW INSTANTIATING RUNTIME ENVIRONMENT

# Target container
FROM scratch

LABEL maintainer="ivann@ntnu.no"

# Root as working directory to copy compiled file to
WORKDIR /

# Retrieve binary from builder container
COPY --from=builder /go/src/app/cmd/server .

EXPOSE 8080

# Instantiate binary
CMD ["./server"]
