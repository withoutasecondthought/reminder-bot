package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("error from env load")
	}

	viper.AutomaticEnv()

	bot, err := tgbotapi.NewBotAPI(viper.GetString("TOKEN"))
	if err != nil {
		log.Fatal().Err(err).Msg("error from new bot api")
	}

	// uncomment for debug
	//bot.Debug = true

	log.Info().Msg("Authorized on account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, viper.GetString("HELP_MESSAGE"))
			msg.ReplyToMessageID = update.Message.MessageID

			if _, err := bot.Send(msg); err != nil {
				log.Error().Err(err).Msg("error from help")
			}
		}

		if update.Message.Command() == "remind" {
			data := strings.Split(update.Message.CommandArguments(), " ")
			log.Info().Msg(fmt.Sprintf("data: %v", data))

			if len(data) < 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, viper.GetString("BAD_REQUEST_MESSAGE"))
				msg.ReplyToMessageID = update.Message.MessageID

				if _, err := bot.Send(msg); err != nil {
					log.Error().Err(err).Msg("error from remind bad request")
				}

				continue
			}

			duration, err := time.ParseDuration(data[0])
			if err != nil {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					fmt.Sprintf("%s: %s", viper.GetString("WRONG_TIME_MESSAGE"), err),
				)
				msg.ReplyToMessageID = update.Message.MessageID

				if _, err := bot.Send(msg); err != nil {
					log.Error().Err(err).Msg("error from remind parse duration")
				}

				continue
			}

			var remindNicknames = make(map[string]struct{})

			var originFlag, meFlag = false, false

			for _, user := range data[1:] {
				if user == "me" && !meFlag {
					remindNicknames["@"+update.Message.From.UserName] = struct{}{}
					meFlag = true

					continue
				}

				if user == "origin" && !originFlag {
					for _, nickname := range findNickNames(update.Message.ReplyToMessage.Text) {
						remindNicknames[nickname] = struct{}{}
					}

					originFlag = true

					continue
				}

				remindNicknames[user] = struct{}{}
			}

			var uniqueNicknames []string

			for nickname := range remindNicknames {
				uniqueNicknames = append(uniqueNicknames, nickname)
			}

			prepareRemindFunc := func() {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					fmt.Sprintf("%s %s", viper.GetString("REMIND_FOR_MESSAGE"), strings.Join(uniqueNicknames, " ")),
				)

				if update.Message.ReplyToMessage != nil {
					msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
				} else {
					msg.ReplyToMessageID = update.Message.MessageID
				}

				time.Sleep(duration)

				if _, err := bot.Send(msg); err != nil {
					log.Error().Err(err).Msg("error from remind send remind")
				}
			}

			go prepareRemindFunc()

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("%s %s", viper.GetString("REMIND_IN_MESSAGE"), duration),
			)
			msg.ReplyToMessageID = update.Message.MessageID

			if _, err = bot.Send(msg); err != nil {
				log.Error().Err(err).Msg("error from remind send success")
			}
		}
	}
}

func findNickNames(txt string) []string {
	var nicknames []string

	for _, word := range strings.Split(txt, " ") {
		if strings.HasPrefix(word, "@") {
			nicknames = append(nicknames, word)
		}
	}

	return nicknames
}
