FROM golang:latest

COPY ingester.go ./
RUN go build -o /usr/bin/ingester ingester.go

CMD ["/usr/bin/ingester"]
