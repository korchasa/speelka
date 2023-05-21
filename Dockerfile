FROM golang:1.20.0-alpine3.17 as build
COPY --from=golangci/golangci-lint:v1.51.1 /usr/bin/golangci-lint /usr/bin/golangci-lint
WORKDIR /build
ENV \
    TERM=xterm-color \
    TIME_ZONE="UTC" \
    CGO_ENABLED=0 \
    GOFLAGS="-mod=vendor"

COPY go.* ./
COPY bin/app/main.go ./main.go
COPY vendor/ ./vendor
COPY pkg ./pkg

RUN go env
RUN go version
RUN echo "  ## Test" && go test ./...
RUN echo "  ## Lint" && golangci-lint run ./...
RUN echo "  ## Build" && go build -o app .

FROM alpine:3.15
WORKDIR /app
COPY --from=build /build/app ./app
USER nobody:nobody
CMD ["./app"]