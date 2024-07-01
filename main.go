package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JuanJDlp/File_Storage_System/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := api.NewRouter()

	log.Fatal(router.Start(os.Getenv("PORT")))

	fmt.Print("Server Started on port")
}
