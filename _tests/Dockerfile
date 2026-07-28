FROM golang:1.22-alpine
WORKDIR /app
RUN apk add --no-cache curl >/dev/null 2>&1
COPY . .
RUN go build -o /usr/local/bin/dreego ./cmd/dreego
ENTRYPOINT ["sh", "_tests/test.sh"]
