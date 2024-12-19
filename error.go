package main

import (
	"fmt"
	"time"
)

func processError(err error) {
	if err != nil {
		logger.Error().Err(err).Msg("Unexpected Error")
		time.Sleep(time.Second * 5)
		panic(err)
	}
}

func processErrorWithoutLogger(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		time.Sleep(time.Second * 5)
		panic(err)
	}
}
