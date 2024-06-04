FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o myapp ./cmd/main.go

ENV DEV=true
ENV PORT=8080

EXPOSE $PORT

ENV FUNCTION_TARGET=WhatIsMyIP

CMD ["./myapp"]
