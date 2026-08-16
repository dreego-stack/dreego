FROM golang:1.22-alpine
ARG DREEGO_VERSION=dev
WORKDIR /app
RUN apk add --no-cache curl py3-yaml >/dev/null 2>&1
COPY . .
RUN go build -ldflags "-X main.version=$DREEGO_VERSION" -o /usr/local/bin/dreego ./cli/dreego
ENTRYPOINT ["sh", "_tests/test.sh"]
