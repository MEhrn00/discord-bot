FROM golang:1.15-alpine

RUN mkdir -p /usr/src/app

RUN apk add git

WORKDIR /usr/src/app
RUN mkdir commands

COPY bot.go .
COPY commands/commands.go commands
COPY go.mod .
COPY go.sum .

RUN go build -o discord-bot

CMD ["./discord-bot"]
