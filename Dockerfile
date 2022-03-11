FROM golang:1.12-alpine as build

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/kainlite/whatismyip-go
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

ENV \
  CGO_ENABLED=0 \
  GOOS=linux

# Install the package and create test binary
RUN go install -v ./... && \
    go test -c


FROM scratch
COPY --from=build /go/bin/whatismyip-go /whatismyip

# This container exposes port 8080 to the outside world
EXPOSE 8000

USER 1000

# Run the executable
CMD ["/whatismyip"]
