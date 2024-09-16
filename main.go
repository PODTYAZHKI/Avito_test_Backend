package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/joho/godotenv"

	"tender-service/app"
	"tender-service/database"
	"tender-service/routes"
)


func main() {
	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Подключение к базе данных
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	} else {
		log.Println("Database connected!")
	}
	defer db.Close()

	// Выполнение миграций
	database.Migrate(db)
	log.Println("Database migrated successfully!")
	myApp := app.NewApp(db)

	defaultRouter := routes.NewDefaultRouter(myApp)


	serverAddress := os.Getenv("SERVER_ADDRESS")
	if err := http.ListenAndServe(serverAddress, defaultRouter); err != nil {
		log.Fatalf("Could not start server: %v", err)
	} else {
		log.Println("Server is running on", os.Getenv("SERVER_ADDRESS"))
	}

}
