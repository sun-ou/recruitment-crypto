# Compile stage
FROM golang:1.23-alpine3.20 AS builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk add --no-cache git make bash ca-certificates tzdata \
    --repository http://mirrors.aliyun.com/alpine/v3.11/community \
    --repository http://mirrors.aliyun.com/alpine/v3.11/main

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.8 && \
    wget -qO/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct" \
    TZ=Asia/Shanghai \
    APP_ENV=docker

WORKDIR /go/src/github.com/sun-ou/

# Copy go.mod, go.sum, download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy project files
COPY . .

# Build the Go app
RUN go build -o recruitment-crypto

# Create minimal image
# Final stage
FROM alpine:3.20
 
WORKDIR /bin
ENV APP_ENV local

COPY --from=builder /go/src/github.com/sun-ou/recruitment-crypto    /bin/recruitment-crypto
COPY --from=builder /bin/grpc_health_probe                          /bin/grpc_health_probe

RUN apk update \
 && apk add --no-cache curl jq \
 && rm -rf /var/cache/apk/* \
 && mkdir -p  /data/logs/

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the app
CMD ["/bin/recruitment-crypto"]

# 1. build image: docker build -t candidate-sun:v1.0.1 -f Dockerfile .
# 2. start: docker run --rm -it -p 8080:8080 candidate-sun:v1
# 3. test: curl -i http://localhost:8080/health
