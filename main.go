package main

import (
	"fmt"
	"os"

	"github.com/glebsnigirev/final-GS/pkg/db"
	"github.com/glebsnigirev/final-GS/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		fmt.Println("Ошибка при инициализации базы данных:", err)
		return
	}

	err = server.Run()
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
}
