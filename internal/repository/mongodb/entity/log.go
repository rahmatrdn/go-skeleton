package entity

import "time"

type LogCollection struct {
	Status        string            `bson:"status" json:"status"`
	Message       string            `bson:"message" json:"message"`
	FuncName      string            `bson:"func_name" json:"func_name"`
	ErrorMessage  string            `bson:"error_message" json:"error_message"`
	Process       string            `bson:"process_name" json:"process_name"`
	LogFields     map[string]string `bson:"log_fields" json:"log_fields"`
	Created       time.Time         `bson:"created" json:"created"`
	ExecutionTime int               `bson:"exec_time" json:"exec_time"`
}

func NewLogCollection() LogCollection {
	instance := LogCollection{}
	return instance
}
