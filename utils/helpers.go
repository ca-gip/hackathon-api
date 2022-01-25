package utils

import (
	"github.com/rs/zerolog/log"
	_ "os"
)

func Checkb(b bool, msg string) {
	if !b {
		log.Error().Msgf("%v ", msg)
	}
}
