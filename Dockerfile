FROM golang:1.20-alpine3.19 AS builder

COPY . /github.com/SadGodSee/telegram-bot-one/
WORKDIR /github.com/SadGodSee/telegram-bot-one

RUN go mod download
RUN go build -o ./bin/ cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/SadGodSee/telegram-bot-one/bin/bot .
COPY --from=0 /github.com/SadGodSee/telegram-bot-one/configs configs/

EXPOSE 80

CMD ["./bot"]