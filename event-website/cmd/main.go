package main

import (
	"log"
	"net/http"

	"event-website/internal/handlers"
)

func main() {
	http.HandleFunc("/home", handlers.IndexHandler)
	http.HandleFunc("/create", handlers.CreateEventHandler)
	http.HandleFunc("/edit/", handlers.EditEventHandler)
	http.HandleFunc("/delete/", handlers.DeleteEventHandler) //Доделать функцию удаления 6.01
	http.HandleFunc("/", handlers.HomeHandler)               // 5.01 Регистр сделан
	http.HandleFunc("/view_page", handlers.ViewPageHandler)  // 10.01 просмотр

	http.HandleFunc("/sign_in", handlers.SignHandler) // 12.01 просмотр
	http.HandleFunc("/log_in", handlers.LogHandler)

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	log.Println("Сервер запущен. Откройте http://localhost:8080 в вашем браузере.")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Дописать ридми 11.01
//Доделать запуск через екзешник 11.01
