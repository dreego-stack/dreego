FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /dreego ./cli/dreego
RUN cd demo && /dreego generate && CGO_ENABLED=0 go build -o /app -ldflags="-s -w" .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
