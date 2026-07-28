FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init §$name$§ && go mod tidy && CGO_ENABLED=0 go build -o /app/bin/server .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/bin/server /server
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/server"]
