package helper

import (
	"log"
	"os"
)

func WriteLog(data string, channel string) error {
	f, err := os.OpenFile(channel,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		log.Println(err)
	}

	return err
}
