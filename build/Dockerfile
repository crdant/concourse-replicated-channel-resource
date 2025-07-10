FROM cgr.dev/chainguard/go:latest-dev AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o concourse-replicated-channel-resource ./cmd/concourse-replicated-channel-resource

FROM cgr.dev/chainguard/static:latest

COPY --from=builder /app/concourse-replicated-channel-resource /opt/resource/check
COPY --from=builder /app/concourse-replicated-channel-resource /opt/resource/in
COPY --from=builder /app/concourse-replicated-channel-resource /opt/resource/out

USER 65532:65532