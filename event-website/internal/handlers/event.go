package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"event-website/internal/models"
	"event-website/internal/storage"
)

var storageInstance = storage.NewEventStorage()

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var users []User

const usersFile = "users.json"

func init() {
	if err := storageInstance.Load("events.json"); err != nil {
		log.Fatalf("Ошибка загрузки событий: %v", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на index")
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	events := storageInstance.GetAll()

	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		log.Printf("Поиск событий по запросу: %s", searchQuery)
		filteredEvents := []models.Event{}
		for _, event := range events {
			if strings.Contains(strings.ToLower(event.Title), strings.ToLower(searchQuery)) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}

	if err := tmpl.Execute(w, events); err != nil {
		log.Printf("Ошибка отображения страницы: %v", err)
		http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
	}
}

// func contains(s, substr string) bool {
// 	return strings.Contains(strings.ToLower(s), strings.ToLower(substr)) //Дороботать 12.01 ненужное но и нужное
// }

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		category := r.FormValue("category")

		if title == "" || description == "" || category == "" {
			log.Println("Ошибка: поля не могут быть пустыми")
			http.Error(w, "Поля не могут быть пустыми", http.StatusBadRequest)
			return
		}

		event := models.NewEvent(title, description, category)
		storageInstance.Add(event)

		if err := storageInstance.Save("events.json"); err != nil {
			log.Printf("Ошибка сохранения события: %v", err)
			http.Error(w, "Ошибка сохранения события", http.StatusInternalServerError)
			return
		}
		log.Println("Событие создано:", event.Title)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/create_event.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Ошибка отображения create_event: %v", err)
		http.Error(w, "Ошибка отображения create_event", http.StatusInternalServerError)
	}
}

func EditEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/"):]

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		category := r.FormValue("category")

		if title == "" || description == "" || category == "" {
			http.Error(w, "Поля не могут быть пустыми", http.StatusBadRequest)
			return
		}

		event := models.Event{
			ID:          id,
			Title:       title,
			Description: description,
			Category:    category,
		}
		storageInstance.Edit(event)

		if err := storageInstance.Save("events.json"); err != nil {
			http.Error(w, "Ошибка сохранения события", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/edit_event.html"))
	event, exists := storageInstance.GetEventByID(id)
	if exists {
		if err := tmpl.Execute(w, event); err != nil {
			http.Error(w, "Ошибка отображения edit_event", http.StatusInternalServerError)
		}
		return
	}
	http.NotFound(w, r)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/delete/"):]

	if storageInstance.Delete(id) {
		if err := storageInstance.Save("events.json"); err != nil {
			http.Error(w, "Ошибка сохранения изменений", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Событие не найдено", http.StatusNotFound)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на home")
	tmpl := template.Must(template.ParseFiles("web/templates/home.html"))
	events := storageInstance.GetAll()
	if err := tmpl.Execute(w, events); err != nil {
		log.Printf("Ошибка отображения home: %v", err)
		http.Error(w, "Ошибка отображения home", http.StatusInternalServerError)
	}
}

func ViewPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на view_page")
	tmpl := template.Must(template.ParseFiles("web/templates/view_page.html"))
	events := storageInstance.GetAll()
	if err := tmpl.Execute(w, events); err != nil {
		log.Printf("Ошибка отображения view_page: %v", err)
		http.Error(w, "Ошибка отображения view_page", http.StatusInternalServerError)
	}
}

func LoadUsers() error {
	file, err := ioutil.ReadFile(usersFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(file, &users)
}

func SaveUsers() error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(usersFile, data, 0644)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")

		for _, user := range users {
			if user.Username == username {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		newUser := User{Username: username, Password: password, Role: role}
		users = append(users, newUser)

		if err := SaveUsers(); err != nil {
			http.Error(w, "Ошибка при сохранении пользователя", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/log_in", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		for _, user := range users {
			if user.Username == username && user.Password == password {
				if user.Role == "student" {
					http.Redirect(w, r, "/view_page", http.StatusSeeOther)
				} else if user.Role == "teacher" {
					http.Redirect(w, r, "/home", http.StatusSeeOther)
				}
				return
			}
		}
		http.Redirect(w, r, "/log_in", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/log_in", http.StatusSeeOther)

}

func SignHandler(w http.ResponseWriter, r *http.Request) {
	if err := LoadUsers(); err != nil {
		fmt.Println("Ошибка при загрузке пользователей:", err) //11.01 Полня регмистрация
		return
	}

	log.Println("Получен запрос на sign_in")
	tmpl := template.Must(template.ParseFiles("web/templates/sign_in.html"))
	events := storageInstance.GetAll()
	if err := tmpl.Execute(w, events); err != nil {
		log.Printf("Ошибка отображения sign_in: %v", err)
		http.Error(w, "Ошибка отображения sign_in", http.StatusInternalServerError)
	}
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	if err := LoadUsers(); err != nil {
		fmt.Println("Ошибка при загрузке пользователей:", err) //11.01 Полня регмистрация
		return
	}

	log.Println("Получен запрос на log_in")
	tmpl := template.Must(template.ParseFiles("web/templates/log_in.html"))
	events := storageInstance.GetAll()
	if err := tmpl.Execute(w, events); err != nil {
		log.Printf("Ошибка отображения log_in: %v", err)
		http.Error(w, "Ошибка отображения log_in", http.StatusInternalServerError)
	}
}
