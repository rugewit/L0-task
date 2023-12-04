package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"subscriber/additional"
	"subscriber/api"
	"subscriber/db"
	"subscriber/services"
	"subscriber/subscriber"
	"time"
)

func main() {
	fmt.Println("Hello backend!")

	err := additional.LoadViper("../env/.env")
	if err != nil {
		log.Fatalln("cannot load viper")
		return
	}

	apiPort := viper.Get("API_PORT").(string)
	if err != nil {
		log.Fatal(err)
		return
	}
	// init db
	ctx := context.Background()
	postgres, err := db.NewPostgresDb(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer postgres.Close(ctx)

	// init order service
	orderService := services.NewPostgresOrderService(postgres)

	// init cache
	defaultExpiration, err := additional.GetIntVariableFromViper("DEFAULT_CACHE_TIME_MINUTES")
	if err != nil {
		log.Fatalln("cannot get DEFAULT_CACHE_TIME_MINUTES")
		return
	}
	cleanupInterval, err := additional.GetIntVariableFromViper("DEFAULT_CACHE_CLEANUP_MINUTES")
	if err != nil {
		log.Fatalln("cannot get DEFAULT_CACHE_CLEANUP_MINUTES")
		return
	}

	cacheService := services.NewRWMCache(time.Duration(defaultExpiration)*time.Minute,
		time.Duration(cleanupInterval)*time.Minute)
	// try to restore cache
	restoreCount, err := additional.GetIntVariableFromViper("ORDERS_CACHE_RESTORE_COUNT")
	if err != nil {
		log.Fatalln("cannot get ORDERS_CACHE_RESTORE_COUNT")
		return
	}
	cacheService.Restore(orderService, restoreCount)

	// init subscriber
	subscribeHandler := subscriber.NewSubscribeHandler()
	err = subscribeHandler.SubscribeToMessages(orderService, cacheService)
	if err != nil {
		log.Fatalln("cannot subscribe")
		return
	}

	// init rest api
	router := gin.Default()
	api.RegisterRoutes(router, orderService, cacheService)
	router.Run(apiPort)
}
