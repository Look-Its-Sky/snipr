package main

import (
	"log"
	"os"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var db *gorm.DB 

func init() {
	// Initial connection
	host := os.Getenv("DB_HOST")  
	port := os.Getenv("DB_PORT") 
	user := os.Getenv("DB_USER") 
	password := os.Getenv("DB_PASS") 
	name := os.Getenv("DB_NAME") 

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database!\n %v", err)
	}

	err = db.AutoMigrate(&Contract{})
	if err != nil {
		log.Fatalf("Failed to migrate database!\n %v", err)
	}

	log.Println("Connection to database was successful!")
}

func pushNew(c Contract) {
	result := db.Create(&c)
	if result.Error != nil {
		log.Printf("Error pushing contract to db: %v", result.Error)
	} else {
		log.Printf("Pushed %s to db", c.ContractAddress)
	}

	// fmt.Printf("%#v\n", c)
	// panic("STOP")
}
