FROM rust:latest

WORKDIR /opt/agent
COPY Cargo.toml Cargo.lock ./
COPY src/ ./src
RUN cargo build --release
RUN mv target/release/agent /usr/bin/agent

CMD ["/usr/bin/agent"]
