# Stage 1: Build stage
FROM golang:1.22-alpine3.20 AS build

# Set the working directory
WORKDIR /server

# Copy and download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o webserver cmd/webserver/webserver.go

# Stage 2: Final stage
FROM alpine:latest AS build-release-stage

# Set the working directory
WORKDIR /server

# Copy the binary from the build stage
COPY --from=build /server/webserver .

# Set the entrypoint command
ENTRYPOINT ["/server/webserver"]