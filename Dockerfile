FROM golang:1.19 as builder
WORKDIR /app
COPY . .
WORKDIR /app/bin
RUN go build -o server cmd/http/main.go

FROM scratch
COPY --from=builder /app/bin/server /server
ENTRYPOINT ["/server"]
