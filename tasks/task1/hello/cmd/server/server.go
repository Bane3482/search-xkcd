package main

import (
	"flag"
	"fmt"
	"net/http"

	"yadro.com/course/config"
	"yadro.com/course/utils"
)

func main() {
	configPath := flag.String("config", "", "config for setting app")

	cfg, err := config.New(*configPath)

	if err != nil {
		fmt.Println("error config")
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/ping", &utils.HelloHandler{})
	mux.Handle("/hello", &utils.HelloHandler{})

	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		return
	}
}
