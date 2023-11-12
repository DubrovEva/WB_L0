package entities

import (
	"net/http"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"` // TODO: проверить, что дата будет корректно обрабатываться
	OofShard          string    `json:"oof_shard"`
	Delivery          `json:"delivery"`
	Payment           `json:"payment"`
	Items             []Item `json:"items"` // TODO: проверить, что списки будут корректно обрабатываться
}

func (o Order) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
