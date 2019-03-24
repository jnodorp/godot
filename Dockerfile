FROM golang:1.12 as builder
LABEL maintainer="Julian Schlichtholz <julian.schlichtholz@gmail.com>"

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

# Copy the sources.
RUN mkdir /app
WORKDIR /app
COPY . .

# Build godot.
RUN go build

FROM scratch
LABEL maintainer="Julian Schlichtholz <julian.schlichtholz@gmail.com>"

COPY --from=builder /app/godot /usr/bin/godot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/usr/bin/godot"]
