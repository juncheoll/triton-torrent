# Use Golang Version
FROM golang:latest

# Working Directory Setting
WORKDIR /torrent

# Base File Copy
COPY . .
RUN go mod download

# Build
RUN go build -o main .

# Execute
CMD ["./main"]