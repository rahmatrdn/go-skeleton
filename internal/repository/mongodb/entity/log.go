package entity

type LogCollection struct {
	Status       string                 `bson:"status" json:"status"`
	FuncName     string                 `bson:"func_name" json:"func_name"`
	ErrorMessage string                 `bson:"error_message" json:"error_message"`
	Process      string                 `bson:"process_name" json:"process_name"`
	LogFields    map[string]interface{} `bson:"log_fields" json:"log_fields"`
}

func NewLogCollection() LogCollection {
	instance := LogCollection{}
	return instance
}
