FROM golang:latest

COPY ingestor.go ./
RUN go build -o /usr/bin/ingestor ingestor.go

CMD ["/usr/bin/ingestor"]
