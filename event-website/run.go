package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "cmd/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	//log.Println("Сервер запущен на http://localhost:8080") ДОбавить и дописать в основу 7.01

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}
}
