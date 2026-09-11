FROM golang:1.27-alpine
ARG DREEGO_VERSION=dev
WORKDIR /app
RUN apk add --no-cache curl git py3-yaml >/dev/null 2>&1
COPY . .
RUN go build -ldflags "-X main.version=$DREEGO_VERSION" -o /usr/local/bin/dreego ./cmd/dreego
ENTRYPOINT ["sh", "_tests/test.sh"]
