package logger

import (
	"errors"
	"log"
	"webtest/internal/models"
)

var SendMsgErr = errors.New("Не удалось отправить сообщение брокеру!")

func LogPrint(message string, level models.LogLevel) error {
	log.Printf("%s\n", message)

	err := NewMessage(models.LogMessage{
		Level:   level,
		Message: message,
	})
	if err != nil {
		return SendMsgErr
	}

	return nil
}
