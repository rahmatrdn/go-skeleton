package helper

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rahmatrdn/go-skeleton/config"
)

func SendNotif(message string) (result string) {
	cfg := config.NewConfig()

	base, err := url.Parse(cfg.TelegramBotOption.Url + "bot" + cfg.TelegramBotOption.Token + "/sendMessage?")
	if err != nil {
		return
	}
	// Query params
	params := url.Values{}
	params.Add("chat_id", cfg.TelegramBotOption.ChatID)
	params.Add("text", message)
	params.Add("parse_mode", "HTML")
	base.RawQuery = params.Encode()
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", base.String(), nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		result = err.Error()
	}

	res, err := client.Do(request)
	if err != nil {
		result = err.Error()
	}
	ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		result = err.Error()
	} else {
		result = string(body)
	}
	return result
}
