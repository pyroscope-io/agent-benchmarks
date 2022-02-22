FROM golang:1.17.7

WORKDIR /opt/agent
COPY main.go go.mod go.sum ./
RUN go build -o /usr/bin/agent
ENV GOMAXPROCS=1

CMD ["/usr/bin/agent"]
