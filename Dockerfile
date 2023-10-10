FROM scratch
COPY nomad-events-sink .
COPY config.sample.toml ./config.toml
CMD ["./nomad-events-sink"]
