package inmemory

type Request struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

