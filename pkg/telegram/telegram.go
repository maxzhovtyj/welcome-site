package telegram

import (
	"errors"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"strings"
	"sync"
)

var (
	tmToken  = flag.String("tmToken", "", "Telegram token to access the HTTP API")
	tmChatID = flag.Int64("tmChatID", 00000000, "Telegram ChatID")
)

var (
	bot             *tgbotapi.BotAPI
	appName         string
	updatesChanOnce sync.Once
	updatesChan     tgbotapi.UpdatesChannel
)

func Init(appNameParam string) error {
	if appNameParam == "" {
		return errors.New("appName is empty")
	}

	var err error
	appName = appNameParam
	bot, err = tgbotapi.NewBotAPI(*tmToken)
	return err
}

func GetUpdatesChan() (tgbotapi.UpdatesChannel, error) {
	var err error
	updatesChanOnce.Do(func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 0
		updatesChan, err = bot.GetUpdatesChan(u)
		if err != nil {
			return
		}
	})
	return updatesChan, err
}

func Send(message string) error {
	if strings.Contains(appName, "___") {
		// it's an internal service
		// don't send message
		return fmt.Errorf("sending of message in the debug mode is forbidden")
	}
	return send(message, nil)
}

func SendRaw(message string) error {
	msg := tgbotapi.NewMessage(*tmChatID, message)
	_, err := bot.Send(msg)
	return err
}

func SendWithLink(message, link string) error {
	uri, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("error while parsing link %q: %s", link, err)
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: uri.Host,
				URL:  &link,
			},
		),
	)
	return send(message, &markup)
}

func send(message string, markup *tgbotapi.InlineKeyboardMarkup) error {
	message = fmt.Sprintf("*%s:* %s", appName, message)
	msg := tgbotapi.NewMessage(*tmChatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	_, err := bot.Send(msg)
	return err
}
