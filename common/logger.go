package common

import (
	"log"
)

type Logger struct {
}

func (l Logger) Error(err error, content ...any) error {
	log.Default().Println(err)
	log.Default().Println(content...)
	return err
}

func (l Logger) Info(content ...any) {
	log.Default().Println(content...)
}