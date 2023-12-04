package subscriber

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"log"
	"subscriber/additional"
	"subscriber/model"
	"subscriber/services/interfaces"
)

const clientId = "subscriber-client"

type SubscribeHandler struct {
	ClusterId string
	ClientId  string
	Channel   string
}

func NewSubscribeHandler() *SubscribeHandler {
	err := additional.LoadViper("../env/.env")
	if err != nil {
		log.Fatalln("cannot load viper")
		return nil
	}

	clusterId := viper.Get("CLUSTER_ID").(string)
	channel := viper.Get("CHANNEL").(string)
	return &SubscribeHandler{
		ClusterId: clusterId,
		ClientId:  clientId,
		Channel:   channel,
	}
}

func WriteToDb(order *model.Order, orderService interfaces.OrderService) error {
	err := orderService.Insert(order, context.Background())
	if err != nil {
		log.Fatalln("cannot insert a model to the db", err)
		return err
	}
	log.Printf("order has been written to the db: %s\n", order.OrderUID)
	return nil
}

func messageCallbackWrapper(orderService interfaces.OrderService, cache interfaces.Cache) func(*stan.Msg) {
	return func(message *stan.Msg) {
		var order model.Order
		err := json.Unmarshal(message.Data, &order)
		if err != nil {
			log.Fatalln("cannot unmarshal message.data into order")
		}
		err = WriteToDb(&order, orderService)
		if err != nil {
			log.Fatalln("cannot write to the db")
		}
	}
}

func (handler *SubscribeHandler) SubscribeToMessages(orderService interfaces.OrderService, cache interfaces.Cache) error {
	natsStreaming, err := stan.Connect(handler.ClusterId, handler.ClientId)
	if err != nil {
		log.Fatalln("cannot connect to nats streaming")
		return err
	}
	_, err = natsStreaming.Subscribe(handler.Channel, messageCallbackWrapper(orderService, cache))
	if err != nil {
		log.Fatalln("cannot subscribe to nats streaming channel ", handler.Channel)
		return err
	}
	log.Println("Subscription has been set")
	return nil
}
