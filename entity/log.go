package entity

import "encoding/json"

const (
	LogGeneralKey = "./storage/log/general.log"
	LogWorkerKey  = "./storage/log/worker.log"
)

type LogType string

const (
	LogSuccess LogType = "SUCCESS"
	LogError   LogType = "ERROR"
	LogInfo    LogType = "INFO"
	LogWarning LogType = "WARNING"
	LogDebug   LogType = "DEBUG"
)

type Log struct {
	FuncName     string            `json:"func_name"`
	Message      string            `json:"message"`
	ErrorMessage string            `json:"error_message"`
	Process      string            `json:"process"`
	Status       LogType           `json:"status"` // ERROR, SUCCESS, WARNING, INFO, DEBUG
	LogFields    map[string]string `json:"log_fields"`
}

func (c *Log) LoadFromMap(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(data, c)
	}
	return err
}
