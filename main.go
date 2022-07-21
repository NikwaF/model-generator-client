package main

import (
	"github.com/joho/godotenv"
	"echo-api/routes"
)

func main() {
	godotenv.Load()

	web := routes.Init()
	web.Logger.Fatal(web.Start(":8082"))
}

