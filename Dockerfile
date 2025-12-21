# FROM golang:1.24-alpine AS builder
# WORKDIR /app
# RUN apk add --no-cache git
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# FROM alpine:3.19
# RUN apk add --no-cache ca-certificates bash curl wget
# # Install kubectl
# RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
#     chmod +x kubectl && \
#     mv kubectl /usr/local/bin/
# WORKDIR /
# COPY --from=builder /server /server
# COPY --from=builder /app/scripts /scripts
# EXPOSE 8080
# ENTRYPOINT ["/server"]

# ---------- Stage 1: Build ----------
    FROM golang:1.24.10 AS builder
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o dbaas-server ./cmd/server
    # ---------- Stage 2: Runtime ----------
    FROM ubuntu:22.04
    ENV DEBIAN_FRONTEND=noninteractive
    # System deps
    RUN apt-get update && apt-get install -y \
        ca-certificates \
        curl \
        wget \
        bash \
        sshpass \
        openssh-client \
        iproute2 \
        jq \
        git \
        sudo \
        && rm -rf /var/lib/apt/lists/*
    # kubectl
    RUN curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl \
        && chmod +x kubectl \
        && mv kubectl /usr/local/bin/
    # helm
    RUN curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    # velero
    RUN wget -q https://github.com/vmware-tanzu/velero/releases/download/v1.14.0/velero-v1.14.0-linux-amd64.tar.gz \
        && tar -xzf velero-v1.14.0-linux-amd64.tar.gz \
        && chmod +x velero-v1.14.0-linux-amd64/velero \
        && mv velero-v1.14.0-linux-amd64/velero /usr/local/bin/velero \
        && rm -rf velero-v1.14.0-linux-amd64*
    # non-root user
    RUN useradd -m appuser && echo "appuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
    WORKDIR /app
    # binary
    COPY --from=builder /app/dbaas-server /app/dbaas-server
    # runtime files (paths unchanged)
    COPY scripts ./scripts
    COPY kubeconfig.yaml ./kubeconfig.yaml
    COPY credentials-velero ./credentials-velero
    RUN chmod +x ./scripts/*.sh \
        && chown -R appuser:appuser /app
    USER appuser
    EXPOSE 8080
    CMD ["./dbaas-server"]