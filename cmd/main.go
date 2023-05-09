package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aidarbek436/forum/internal/config"
	"github.com/aidarbek436/forum/internal/handler"
	"github.com/aidarbek436/forum/internal/repository"
)

func main() {
	fmt.Println("Server is listening in port localhost:1999")

	cfg, err := config.InitConfig("./config/config.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := repository.OpenDb(cfg.Database)
	if err != nil {
		fmt.Println(err, "with main, opendb")
		return
	}

	defer db.Close()

	storage := repository.NewStorage(db)
	router := handler.NewHandler(storage)

	if err := http.ListenAndServe(":1999", router.Router()); err != nil {
		log.Fatal(err)
	}
}
