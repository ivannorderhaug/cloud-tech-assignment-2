FROM golang:1.17-alpine as Builder

LABEL stage=builder

# Set up execution environment in container's GOPATH
WORKDIR /go/src/app/cmd

# Copy relevant folders into container
COPY ./corona-information-service/cmd /go/src/app/cmd
COPY ./corona-information-service/internal /go/src/app/internal
COPY ./corona-information-service/pkg /go/src/app/pkg
COPY ./corona-information-service/tools /go/src/app/tools
COPY ./corona-information-service/go.mod /go/src/app/go.mod
COPY ./corona-information-service/go.sum /go/src/app/go.sum

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

# To get the time zone data
FROM alpine:latest as alpine-with-tz
RUN apk --no-cache add tzdata zip
WORKDIR /usr/share/zoneinfo

#Compressing the zone data
RUN zip -q -r -0 /zoneinfo.zip .

# Final container
FROM scratch AS final

LABEL maintainer="ivann@ntnu.no"

# Root as working directory to copy compiled file to
WORKDIR /

# Retrieve binary from builder container
COPY --from=builder /go/src/app/cmd/server .
COPY ./service-account.json .

# Setting time zone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine-with-tz /zoneinfo.zip /
ENV TZ=Europe/Berlin

# Fetching the cert hints.
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

# Instantiate server
CMD ["./server"]
