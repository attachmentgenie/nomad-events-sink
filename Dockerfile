FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3.18.4
RUN mkdir -p ./data/events
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY nomad-events-sink .
COPY config.sample.toml ./config.toml
CMD ["./nomad-events-sink"]
