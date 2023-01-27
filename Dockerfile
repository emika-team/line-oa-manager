FROM golang:1.19 as builder

ENV GO111MODULE=on CGO_ENABLED=0

WORKDIR /app
COPY . .
RUN apt update && apt install ca-certificates && update-ca-certificates
RUN go build -o bin/server cmd/http/main.go

FROM scratch
COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/server"]
