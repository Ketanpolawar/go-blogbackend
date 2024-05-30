package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/ketanpolawar/blogbackend/database"
	"github.com/ketanpolawar/blogbackend/routes"
)

func main() {
	database.Connect()
	err:=godotenv.Load()
	if err != nil{
		log.Fatal("Error loading the env file")
	}
	port:=os.Getenv("PORT")
	app:=fiber.New()
	routes.Setup(app)
	app.Listen(":"+port)
}