FROM golang:1.19 as builder

RUN mkdir ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts

ARG PRIVATE_KEY
RUN echo "$PRIVATE_KEY" > ~/.ssh/id_rsa && \
    chmod 400 ~/.ssh/id_rsa

ENV GO111MODULE=on CGO_ENABLED=0

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /app
COPY . .
RUN apt update && apt install ca-certificates && update-ca-certificates
RUN go build -o bin/server cmd/http/main.go

RUN rm -rf ~/.ssh/

FROM scratch
COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/server"]
