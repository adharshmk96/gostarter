package common

type GeneralResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}
