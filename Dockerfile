FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /bin/inference-stub ./cmd/inference-stub

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/inference-stub /usr/local/bin/inference-stub

RUN adduser -D -g '' stubuser
USER stubuser

ENTRYPOINT ["/usr/local/bin/inference-stub"]
