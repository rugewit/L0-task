package main

import (
	"encoding/json"
	"log"
	"publisher/publisher"
	"subscriber/additional"
	"subscriber/model"
	"time"
)

func main() {
	err := additional.LoadViper("../env/.env")
	if err != nil {
		log.Fatalln("cannot load env:", err)
		return
	}

	publishDelay, err := additional.GetIntVariableFromViper("MESSAGE_PUBLISH_DELAY_SEC")
	if err != nil {
		log.Fatalln("cannot get ORDERS_TOTAL_COUNT")
		return
	}

	fileContent, err := publisher.GetJsonFileContent()
	if err != nil {
		log.Fatalln("Error publishing message:", err)
		return
	}
	var order model.Order

	err = json.Unmarshal(fileContent, &order)
	if err != nil {
		log.Fatalln("cannot unmarshal to order")
		return
	}

	publishHandler, err := publisher.NewPublishHandler()
	if err != nil {
		log.Fatalln("cannot create publishHandler", err)
		return
	}
	defer publishHandler.CloseConnection()

	totalCount, err := additional.GetIntVariableFromViper("ORDERS_TOTAL_COUNT")
	if err != nil {
		log.Fatalln("cannot get ORDERS_TOTAL_COUNT")
		return
	}
	publishHandler.PublishMessagesWithDelay(&order, time.Duration(publishDelay)*time.Second, totalCount)
}
