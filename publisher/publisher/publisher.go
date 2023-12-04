package publisher

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"subscriber/additional"
	"subscriber/model"
	"time"
)

const clientId = "publisher-client"

type PublishHandler struct {
	ClusterId     string
	ClientId      string
	Channel       string
	NatsStreaming stan.Conn
}

func NewPublishHandler() (*PublishHandler, error) {
	err := additional.LoadViper("../env/.env")
	if err != nil {
		log.Fatalln("cannot load viper")
		return nil, err
	}

	clusterId := viper.Get("CLUSTER_ID").(string)
	channel := viper.Get("CHANNEL").(string)

	natsStreaming, err := stan.Connect(clusterId, clientId)
	if err != nil {
		log.Fatalln("cannot connect to nats streaming", err)
		return nil, err
	}

	return &PublishHandler{
		ClusterId:     clusterId,
		ClientId:      clientId,
		Channel:       channel,
		NatsStreaming: natsStreaming,
	}, nil
}

func GetJsonFileContent() ([]byte, error) {
	curWorkingPath, err := os.Getwd()
	if err != nil {
		log.Fatalln("cannot get current working directory", err)
		return nil, err
	}
	rootPath := filepath.Dir(curWorkingPath)
	jsonFilePath := filepath.Join(rootPath, "model.json")
	if _, err := os.Stat(jsonFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("json file does not exist!", err)
		return nil, err
	}
	var fileContent []byte
	fileContent, err = os.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatalln("cannot get content of json file", err)
		return nil, err
	}

	return fileContent, nil
}

func (handler *PublishHandler) PublishMessage(message []byte) error {
	return handler.NatsStreaming.Publish(handler.Channel, message)
}

func (handler *PublishHandler) PublishMessagesWithDelay(message *model.Order, delay time.Duration, totalCount int) {
	curCount := 0
	for {
		message.OrderUID = uuid.New().String()
		messageJson, err := json.Marshal(message)
		if err != nil {
			log.Fatalln("cannot convert into json")
			return
		}
		err = handler.PublishMessage(messageJson)
		if err != nil {
			log.Fatalln("cannot publish message")
			return
		}
		log.Println("The message has been sent", message.OrderUID)

		time.Sleep(delay)
		curCount++
		if curCount == totalCount {
			break
		}
	}
}

func (handler *PublishHandler) PublishMessagesWithoutDelay(message *model.Order, totalCount int) {
	curCount := 0
	for {
		message.OrderUID = uuid.New().String()
		messageJson, err := json.Marshal(message)
		if err != nil {
			log.Fatalln("cannot convert into json")
			return
		}

		err = handler.PublishMessage(messageJson)
		if err != nil {
			log.Fatalln("cannot publish message")
			return
		}
		//log.Println("The message has been sent", message.OrderUID)
		curCount++
		if curCount == totalCount {
			break
		}
	}
	log.Println("Messages have been sent")
}

func (handler *PublishHandler) CloseConnection() {
	handler.NatsStreaming.Close()
}
