FROM golang:bullseye as builder

WORKDIR /reminder-bot
COPY . .

RUN CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -v -o ./application ./cmd/*.go

FROM alpine:3.15.4
WORKDIR /reminder-bot

COPY --from=builder /reminder-bot/application /reminder-bot/application
COPY .env .env
CMD ["/reminder-bot/application"]

