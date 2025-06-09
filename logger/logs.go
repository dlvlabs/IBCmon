package logger

import (
	"fmt"

	"github.com/dlvlabs/ibcmon/alert"
	"github.com/rs/zerolog/log"
)

func Info(msg string) {
	log.Info().Msg(msg)
}

func Warn(msg any) {
	message := fmt.Sprint(msg)
	log.Warn().Msg(message)
}

func Error(err error) {
	// send error msg with stack trace to telegram
	alert.SendTg(fmt.Sprintf("%+v", err))

	log.Error().Stack().Err(err).Msg("")
}

func Debug(msg any) {
	message := fmt.Sprint(msg)
	log.Debug().Msg(message)
}
