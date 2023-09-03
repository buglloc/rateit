FROM golang:1.21.0 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/rateit ./cmd/rateit

FROM debian:bookworm-backports

COPY --from=build /go/bin/rateit /usr/sbin/local
ADD example/config.yaml /etc/rateit/config.yaml

ENTRYPOINT ["/usr/sbin/local/rateit --config /etc/rateit/config.yaml"]
