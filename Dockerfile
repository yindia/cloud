# syntax=docker/dockerfile:1.2
FROM cgr.dev/chainguard/go as build

WORKDIR /work

# Use build args for cache keys
ARG CACHEBUST=1

# Copy only go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the rest of the application code
COPY . ./

# Build CLI
FROM build as cli
RUN go build -o /usr/local/bin/task-cli ./cli/main.go

# Build CLI
FROM cgr.dev/chainguard/go  as river
RUN go install github.com/riverqueue/river/cmd/river@latest


# Build Server
FROM build as server
RUN go build -o /usr/local/bin/task-server ./server/root/main.go

# Final image for CLI
FROM cgr.dev/chainguard/go as cli-final
COPY --from=cli /usr/local/bin/task-cli /usr/local/bin/task-cli
ENTRYPOINT ["task-cli"]

# Final image for Server
FROM cgr.dev/chainguard/go as server-final
COPY --from=server /usr/local/bin/task-server /usr/local/bin/task-server
EXPOSE 8080
ENTRYPOINT ["task-server"]


FROM cgr.dev/chainguard/go as river-final
COPY --from=river /root/go/bin/river /usr/local/bin/river
ENTRYPOINT ["river"]