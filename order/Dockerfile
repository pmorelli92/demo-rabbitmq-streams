FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS compiler
RUN apk --update add ca-certificates git

WORKDIR /app
ADD . .

ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /bin/goapp ./cmd/main.go

FROM scratch
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=compiler /bin/goapp /app
USER 65534:65534
ENTRYPOINT ["/app"]
