FROM golang:1.21

WORKDIR /app

COPY topics/create_topics/ .

RUN go mod init create-topics && \
    go get github.com/segmentio/kafka-go gopkg.in/yaml.v3 && \
    go build -o create-topics .

CMD ["./create-topics"]
