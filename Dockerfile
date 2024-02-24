# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

RUN go mod download

# Install the tokenizers library
ARG VERSION=v0.7.0
RUN curl -fsSL https://github.com/daulet/tokenizers/releases/download/${VERSION}/libtokenizers.linux-amd64.tar.gz | tar xvz
RUN mv ./libtokenizers.a /go/pkg/mod/github.com/daulet/tokenizers@${VERSION}/libtokenizers.a


# Build the Go application
RUN go build -o main .

# # Command to run the application
CMD ["./main"]
