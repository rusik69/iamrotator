FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/iamrotator /app/iamrotator
COPY --from=builder /app/front/index.html /app/front/index.html
RUN chmod +x /app/iamrotator
CMD ["/app/iamrotator", "web", "--config", "/app/config.yaml"]