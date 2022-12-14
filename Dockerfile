FROM golang:1.19 as builder
WORKDIR /app
COPY . .
RUN apt update && apt install ca-certificates && update-ca-certificates
RUN CGO_ENABLED=0 go build -o bin/server cmd/http/main.go

FROM scratch
COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/server"]
