FROM quay.io/theauthgear/golang:1.24.9-noble AS stage1
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM ubuntu:noble
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=stage1 /src/authgear-sms-gateway /usr/local/bin/
COPY ./docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 8091
CMD ["authgear-sms-gateway"]
