package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/glebsnigirev/go_final_project_GS/pkg/api"
)

func Run() error {
	api.Init()

	port := 7540
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		_, err := strconv.Atoi(envPort)
		if err != nil {
			fmt.Println("Error converting TODO_PORT to int:", err)
			return err
		}
	}
	http.Handle("/", http.FileServer(http.Dir("web")))
	fmt.Printf("Starting server on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
