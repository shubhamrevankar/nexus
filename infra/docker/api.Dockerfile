FROM golang:1.22-alpine AS builder

WORKDIR /workspace/services/api
COPY services/api/go.mod ./
COPY services/api ./
RUN go build -o /out/api ./cmd/api

FROM alpine:3.21

RUN addgroup -S nexus && adduser -S nexus -G nexus
USER nexus
COPY --from=builder /out/api /usr/local/bin/api
EXPOSE 8080
ENTRYPOINT ["api"]

