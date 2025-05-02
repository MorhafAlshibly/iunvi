# syntax=docker/dockerfile:1

ARG image_suffix

FROM golang:1.24${image_suffix} AS build_image

ARG service
ARG cgo_enabled
ARG flags

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY gen ./gen
COPY cmd/${service} ./cmd/${service}
COPY internal/${service} ./internal/${service}

# Build
RUN if [ "${cgo_enabled}" = "1" ]; then \
        GOMAXPROCS=200 CGO_ENABLED=${cgo_enabled} GOOS=linux go build -ldflags="-extldflags=-static" -o out.exe ./cmd/${service}; \
    else \
        GOMAXPROCS=200 CGO_ENABLED=${cgo_enabled} GOOS=linux go build -o out.exe ./cmd/${service}; \
    fi

# Final image
FROM alpine AS final_image

ARG port

RUN apk add --no-cache ca-certificates

COPY --from=build_image /app/out.exe ./out.exe

# Expose port
EXPOSE ${port}

ENTRYPOINT ["/out.exe"]