package tglogs

import (
	"flag"
	"fmt"
	"log"
	"wedding/pkg/telegram"
)

var telegramEnabled = flag.Bool("telegramEnabled", false, "Enable telegram")

func Init(appName string) {
	if !*telegramEnabled {
		log.Println("Telegram is not enabled")
		return
	}

	err := telegram.Init(appName)
	if err != nil {
		log.Fatalf("cant init telegram. Error: %s", err)
	}

	return
}

type Options struct {
	IsRawStyle bool
}

type OptionsFunc func(o *Options)

func WithRawStyle(isRaw bool) OptionsFunc {
	return func(o *Options) {
		o.IsRawStyle = isRaw
	}
}

func Send(msg string, optsFunc ...OptionsFunc) {
	if !*telegramEnabled {
		return
	}

	var opts Options
	for _, f := range optsFunc {
		f(&opts) //init
	}

	var err error
	if opts.IsRawStyle {
		err = telegram.SendRaw(msg)
	} else {
		err = telegram.Send(msg)
	}

	if err != nil {
		log.Printf("Cant send tm message: %s\n", err)
	}
}

func InitTgBot() {
	if !*telegramEnabled {
		return
	}

	updatesChan, err := telegram.GetUpdatesChan()
	if err != nil {
		msg := fmt.Sprintf("cant get updatesChan. Error: %s", err)
		_ = telegram.Send(msg)
		log.Fatalf(msg)
	}

	for upd := range updatesChan {
		if upd.Message == nil {
			continue
		}

		switch upd.Message.Text {
		case "/ping":
			_ = telegram.Send("pong")
		}
	}
}
