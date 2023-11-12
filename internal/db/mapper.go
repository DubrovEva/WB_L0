package db

import (
	"database/sql"
	"fmt"
	"github.com/DubrovEva/WB_L0/internal/entities"
	"log"
)

func (db Database) SaveOrder(order *entities.Order) error {
	var deliveryID int
	err := db.Conn.QueryRow("INSERT INTO delivery(name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id",
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		return err
	}

	var paymentID int
	err = db.Conn.QueryRow("INSERT INTO payment(transaction, request_id, currency, provider, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING payment_id",
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee).Scan(&paymentID)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard, deliveryID,
		paymentID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = db.Conn.Exec("INSERT INTO item (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUID, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	db.Cache.Store(order.OrderUID, order)

	return nil
}

func (db Database) getOrderFromDB(orderUID string) (*entities.Order, error) {
	var order entities.Order
	var deliveryID, paymentID int
	row := db.Conn.QueryRow("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery, payment FROM orders WHERE order_uid = $1", orderUID)
	err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard, &deliveryID, &paymentID)
	if err != nil {
		return nil, err
	}

	row = db.Conn.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE delivery_id = $1", deliveryID)
	err = row.Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, err
	}

	row = db.Conn.QueryRow("SELECT transaction, request_id, currency, provider, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE payment_id = $1", paymentID)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	rows, err := db.Conn.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1", orderUID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Can't close rows: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return &order, err
}

func (db Database) UpdateCache() error {
	rows, err := db.Conn.Query("SELECT order_uid FROM orders")
	if err != nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Can't close rows: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var orderUID string
		err := rows.Scan(&orderUID)

		if err != nil {
			log.Printf("Can't scan UID: %v\n", err)
		}
		order, err := db.getOrderFromDB(orderUID)
		if err != nil {
			log.Printf("Can't get order with UID %s: %v\n", orderUID, err)
		}

		db.Cache.Store(orderUID, order)
	}

	return nil
}

func (db Database) GetOrderFromCache(orderUID string) (*entities.Order, error) {

	orderAny, ok := db.Cache.Load(orderUID)

	if ok {
		order := orderAny.(*entities.Order)
		log.Printf("Got order from cache")
		return order, nil
	} else {
		orderPtr, err := db.getOrderFromDB(orderUID)
		if err != nil {
			return nil, fmt.Errorf("could not get order with uid: %s", orderUID)
		} else {
			err := db.UpdateCache()
			if err != nil {
				log.Printf("can't update cache: %v", err)
			}
			log.Printf("Got order from db")
			return orderPtr, nil
		}
	}
}
