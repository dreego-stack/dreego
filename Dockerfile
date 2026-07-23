FROM golang:1.22-alpine

WORKDIR /app
COPY go.mod ./
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY demo/ demo/

RUN go build -o /usr/local/bin/dreego ./cmd/dreego
RUN cd demo && dreego generate
RUN go build -o /server ./demo

EXPOSE 8080
ENTRYPOINT ["/server"]
