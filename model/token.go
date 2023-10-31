package model

type Token struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   string `json:"TTL"`
}
