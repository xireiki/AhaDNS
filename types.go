package main

type Answer struct {
	Name string `json:"name"`
	TTL  int    `json:"TTL"`
	Type int    `json:"type"`
	Data string `json:"data"`
}

type Authority Answer
type Additional Answer

type Question struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type DNSEntity struct {
	Status     int          `json:"Status"`
	TC         bool         `json:"TC"`
	RD         bool         `json:"RD"`
	RA         bool         `json:"RA"`
	AD         bool         `json:"AD"`
	CD         bool         `json:"CD"`
	Question   *Question    `json:"Question"`
	Answer     []Answer     `json:"Answer"`
	Authority  []Authority  `json:"Authority"`
	Additional []Additional `json:"Additional"`
	ECS        string       `json:"edns_client_subnet"`
}
