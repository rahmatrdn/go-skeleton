package entity

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
	Message     string `json:"message"`
}

const (
	VALIDATE_CUSTOM_EXAMPLE = "validate_custom_example"
)
