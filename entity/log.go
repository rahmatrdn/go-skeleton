package entity

import "encoding/json"

const (
	GeneralLogFilePath = "/storage/log/general"
	WorkerLogFilePath  = "/storage/log/worker"
)

type Log struct {
	FuncName     string        `json:"func_name"`
	Message      string        `json:"message"`
	ErrorMessage string        `json:"error_message"`
	Process      string        `json:"process"`
	Status       LogType       `json:"status"`
	LogFields    CaptureFields `json:"capture_fields"`
}

type LogType string

const (
	LogSuccess LogType = "SUCCESS"
	LogError   LogType = "ERROR"
	LogInfo    LogType = "INFO"
	LogWarning LogType = "WARNING"
	LogDebug   LogType = "DEBUG"
)

type CaptureFields map[string]string

func (c *Log) LoadFromMap(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(data, c)
	}
	return err
}

const (
	PRODUCTION_ENV = "production"
)
