FROM golang:1.23-alpine

RUN apk add --no-cache \
    ffmpeg \
    yt-dlp \
    opus-dev \
    gcc \
    musl-dev

WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o muzicBot .

# Command to run the application
CMD ["./muzicBot"]