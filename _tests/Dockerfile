FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o /usr/local/bin/dreego ./cmd/dreego
ENTRYPOINT ["sh", "_tests/test.sh"]
