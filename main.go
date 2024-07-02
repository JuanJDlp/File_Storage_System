package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/JuanJDlp/File_Storage_System/api"
	"github.com/joho/godotenv"
)

func main() {
	dbg := flag.Bool("debug", false, "deletes all the files")
	flag.Parse()
	godotenv.Load()

	router := api.NewRouter()

	if *dbg {
		router.Clear()
	}
	log.Fatal(router.Start(os.Getenv("PORT")))

	fmt.Print("Server Started on port")
}
