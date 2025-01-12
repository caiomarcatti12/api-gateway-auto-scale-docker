FROM golang:1.23 AS dev
ENV CONFIG_PATH=/app/configs
WORKDIR /app

FROM dev as build
WORKDIR /app
COPY ./ /app
RUN cd /app/cmd/main && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/api-gateway-auto-scale-docker

FROM alpine:3.18.2 AS prd
ENV CONFIG_PATH=/app/configs
WORKDIR /app
COPY scripts/entrypoint.sh /app/entrypoint.sh
COPY --from=build /app/api-gateway-auto-scale-docker /app/api-gateway-auto-scale-docker

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p /app/configs && \
    chown -R appuser:appgroup /app && \
    chmod +x /app/api-gateway-auto-scale-docker && \
    chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/api-gateway-auto-scale-docker"]
HEALTHCHECK --interval=10s --timeout=3s CMD curl -f http://localhost:8080/ || exit 1