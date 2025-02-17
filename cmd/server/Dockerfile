FROM quay.io/theauthgear/golang:1.23.6-noble AS stage1
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
RUN update-ca-certificates
COPY --from=stage1 /src/authgear-sms-gateway /usr/local/bin/
EXPOSE 8091
CMD ["authgear-sms-gateway"]
