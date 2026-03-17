package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func Start() error {
	portStr := os.Getenv("TODO_PORT")
	var port int

	if portStr == "" {
		port = 7540
	} else {
		var err error
		port, err = strconv.Atoi(portStr)

		if err != nil {
			return fmt.Errorf("неверный формат порта: %s", portStr)
		}
	}

	webDir := "./web"

	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		return fmt.Errorf("папка %s не найдена", webDir)
	}

	fs := http.FileServer(http.Dir(webDir))

	http.Handle("/", fs)

	addr := ":" + strconv.Itoa(port)
	fmt.Printf("Сервер запущен на порту %d\n", port)
	fmt.Printf("Откройте http://localhost:%d в браузере\n", port)

	return http.ListenAndServe(addr, nil)
}
