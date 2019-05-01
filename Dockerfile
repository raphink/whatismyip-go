FROM golang:1.12-alpine

LABEL maintainer="kainlite@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/kainlite/whatismyip-go
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package and create test binary
RUN go install -v ./... && \
    CGO_ENABLED=0 GOOS=linux go test -c

# This container exposes port 8080 to the outside world
EXPOSE 8000

# Set default environment variable values
ENV WHATISMYIP_PORT 8000

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the executable
CMD ["whatismyip-go"]
