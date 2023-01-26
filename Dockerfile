FROM golang:1.19 as builder

ARG PRIVATE_KEY
ARG PUBLIC_KEY

ENV GO111MODULE=on CGO_ENABLED=0

RUN mkdir /root/.ssh

RUN echo "$PRIVATE_KEY" > /root/.ssh/id_rsa && \
    echo "$PUBLIC_KEY" > /root/.ssh/id_rsa.pub && \
    chmod 600 /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa.pub

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /app
COPY . .
RUN apt update && apt install ca-certificates && update-ca-certificates
ARG GITHUB_TOKEN
RUN go build -o bin/server cmd/http/main.go

FROM scratch
COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/server"]
