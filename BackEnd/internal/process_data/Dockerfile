FROM golang:1.20-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY ./process_data/main.go .

# Set necessary environment variables needed for our image and build the process data.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o process_data .

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/process_data", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/process_data"]