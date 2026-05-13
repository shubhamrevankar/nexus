FROM golang:1.23-alpine AS builder

WORKDIR /workspace/services/api
COPY services/api/go.mod services/api/go.sum* ./
COPY services/api ./
RUN go build -o /out/api ./cmd/api

FROM alpine:3.21

RUN addgroup -S nexus && adduser -S nexus -G nexus
WORKDIR /workspace/services/api
USER nexus
COPY --from=builder /out/api /usr/local/bin/api
COPY --from=builder /workspace/services/api/migrations ./migrations
EXPOSE 8080
ENTRYPOINT ["api"]
