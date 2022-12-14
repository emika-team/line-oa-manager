FROM golang:1.19 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/server cmd/http/main.go

FROM scratch
COPY --from=builder /app/bin/server /server
ENTRYPOINT ["/server"]
