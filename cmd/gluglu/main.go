package main

import (
	"log"
	"net/url"
	"os"

	"github.com/baneronetwo/gluglu/internal/ui"
)

func main() {
	// Инициализация логгера
	log.SetOutput(os.Stdout)
	log.SetPrefix("[GluGlu] ")
	log.Println("Запуск браузера GluGlu...")

	// Инициализация UI
	ui, err := ui.New("data:text/html," + url.PathEscape(`
	<html>
		<head><title>gluglu</title></head>
		<body><h1>Hello, world!</h1></body>
	</html>
	`))
	if err != nil {
		log.Fatal(err)
	}
	// defer ui.Close() // This will be added later

	// Запуск UI
	ui.Run()
}