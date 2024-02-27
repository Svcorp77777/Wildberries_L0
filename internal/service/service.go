package service

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"

	"wildberries_L0/internal/model"
	"wildberries_L0/internal/storage"
)

type Service struct {
	storage storage.Storage
	sc      stan.Conn
	cache   map[string][]byte
}

func (ser Service) Subscribe() error {
	_, err := ser.sc.Subscribe("Message", func(m *stan.Msg) {
		var order model.Order

		if ser.validationData(order) {
			err := json.Unmarshal(m.Data, &order)
			if err != nil {
				fmt.Println(err)

				return
			}

			ser.saveCache(order.OrderUid, m.Data)
			err = ser.storage.SaveBD(order.OrderUid, m.Data)
			if err != nil {
				fmt.Println(err)

				return
			}
		}

	})

	if err != nil {
		return err
	}

	return nil
}

func (ser *Service) saveCache(orderUid string, b []byte) {
	ser.cache[orderUid] = b
}

func (ser *Service) GetCache(key string) []byte {
	var answer []byte

	if getCache, exist := ser.cache[key]; exist {
		return getCache
	}

	return answer
}

func (ser *Service) RecoveryCache() error {
	ser.cache = make(map[string][]byte)

	orders, err := ser.storage.RecoveryCacheBD()
	if err != nil {
		return err
	}

	for _, value := range orders {
		ser.cache[value.OrderUid] = value.OrderJson
	}

	return nil
}

func (ser *Service) validationData(data model.Order) (ok bool) {
	if data.OrderUid != "" &&
		data.TrackNumber != "" &&
		data.Entry != "" &&
		data.Delivery.Name != "" &&
		data.Delivery.Phone != "" &&
		data.Delivery.Zip != "" &&
		data.Delivery.City != "" &&
		data.Delivery.Address != "" &&
		data.Delivery.Region != "" &&
		data.Delivery.Email != "" &&
		data.Payment.Transaction != "" &&
		data.Payment.RequestID != "" &&
		data.Payment.Currency != "" &&
		data.Payment.Provider != "" &&
		data.Payment.Amount >= 0 &&
		data.Payment.PaymentDt >= 0 &&
		data.Payment.Bank != "" &&
		data.Payment.DeliveryCost >= 0 &&
		data.Payment.GoodsTotal >= 0 &&
		data.Payment.CustomFee >= 0 &&
		len(data.Items) > 0 &&
		data.Locale != "" &&
		data.InternalSignature == "" &&
		data.CustomerID != "" &&
		data.DeliveryService != "" &&
		data.Shardkey != "" &&
		data.SmID >= 0 &&
		data.DateCreated != "" &&
		data.OofShard != "" {

		for _, value := range data.Items {
			if value.ChrtID >= 0 &&
				value.TrackNumber != "" &&
				value.Price >= 0 &&
				value.Rid != "" &&
				value.Name != "" &&
				value.Sale >= 0 &&
				value.Size != "" &&
				value.TotalPrice >= 0 &&
				value.NmID >= 0 &&
				value.Brand != "" &&
				value.Status >= 0 {

			} else {
				return
			}
		}
	}

	return true
}

func NewService(storage storage.Storage, sc stan.Conn) *Service {
	return &Service{
		storage: storage,
		sc:      sc,
	}
}
