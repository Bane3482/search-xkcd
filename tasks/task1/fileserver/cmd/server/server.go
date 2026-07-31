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

	fmt.Println(cfg.Port)

	mux := http.NewServeMux()

	mux.Handle("/files", &utils.FileHandler{})
	mux.Handle("/files/", &utils.FileHandler{})

	fmt.Println("Star111t")

	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		return
	}
}
