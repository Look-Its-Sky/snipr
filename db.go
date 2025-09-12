package main

import (
	"log"
	"os"
	"fmt"
	"time"

	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/go-redis/redis"
	
	"snipr/schemas"
)

var (
	postgres_db 		*gorm.DB 
	redis_client    *redis.Client
)

func initDB() {
	// Postgres connection
	host := os.Getenv("DB_HOST")  
	port := os.Getenv("DB_PORT") 
	user := os.Getenv("DB_USER") 
	password := os.Getenv("DB_PASS") 
	name := os.Getenv("DB_NAME") 

	postgres_dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	var err error
	postgres_db, err = gorm.Open(postgres.Open(postgres_dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database!\n %v", err)
	}

	err = postgres_db.AutoMigrate(&schemas.Contract{})
	if err != nil {
		log.Fatalf("Failed to migrate database!\n %v", err)
	}

	log.Println("Connection to Postgres was successful!")

	// Redis connection
	redis_addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redis_client = redis.NewClient(&redis.Options{
			Addr:		 redis_addr,
	})

	_, err = redis_client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis!\n %v", err)
	}

	log.Println("Connection to Redis was successful!")
}

func pushNewContract(c *schemas.Contract) {
	// Push to postgres
	result := postgres_db.Create(c)
	if result.Error != nil {
		log.Printf("Error pushing contract to db: %v", result.Error)
	} else {
		if *verbose { log.Printf("Pushed %s to db", c.Address) }
	}

	// Push to Redis
	cacheKey := fmt.Sprintf("%s", c.Address)

  json_data, err := json.Marshal(c)
	if err != nil {
			log.Printf("Error marshaling contract %s to JSON: %v", c.Address, err)
			return
	}

	err = redis_client.Set(cacheKey, json_data, 24 * time.Hour).Err()
	if err != nil {
			log.Printf("Error pushing contract %s to Redis: %v", c.Address, err)
	} else {
			if *verbose { log.Printf("Pushed %s to Redis", c.Address) }
	}
}
