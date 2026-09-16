FROM golang:1.27-alpine
ARG DREEGO_VERSION=dev
WORKDIR /app
RUN apk add --no-cache curl gcc git gtk4.0-dev musl-dev pkgconf py3-yaml strace webkit2gtk-6.0-dev xvfb xwininfo >/dev/null 2>&1
COPY . .
RUN go build -ldflags "-X main.version=$DREEGO_VERSION" -o /usr/local/bin/dreego ./cmd/dreego
ENTRYPOINT ["sh", "_tests/test.sh"]
