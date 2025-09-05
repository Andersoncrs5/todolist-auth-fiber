package res

import "time"

type ResponseHttp[T any] struct {
	Timestamp time.Time `json:"timestamp"`
	Body      T         `json:"body"`
	Code      int       `json:"code"`
	Status    bool      `json:"status"`
	Message   string    `json:"message"`
}
