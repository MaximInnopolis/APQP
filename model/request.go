package model

type EnqueueRequest struct {
	N   int     `json:"n"`
	D   float64 `json:"d"`
	N1  float64 `json:"n1"`
	I   float64 `json:"I"`
	TTL float64 `json:"TTL"`
}
