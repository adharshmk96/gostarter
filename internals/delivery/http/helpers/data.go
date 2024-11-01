package helpers

type GeneralResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type GeneralDataResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
