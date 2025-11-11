package types

type ResponseSingleData[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ResponseMultiData[T any] struct {
	Message string                 `json:"message"`
	Data    []T                    `json:"data"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}
