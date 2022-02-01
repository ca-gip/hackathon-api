package utils

import (
	"errors"
	"github.com/rs/zerolog/log"
	"hackathon-api/models"
	_ "os"
)

func Checkb(b bool, msg string) {
	if !b {
		log.Error().Msgf("%v ", msg)
	}
}

func ValidateMoneyType(moneyType string) error {

	for _, v := range models.GetMoney() {
		if v.MoneyType == moneyType {
			return nil
		}
	}

	return errors.New("Invalid money type")
}
