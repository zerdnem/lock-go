package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/zerdnem/lock-go/utils"
)

var db, _ = gorm.Open("postgres", "")

var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	message := r.PostForm.Get("message")
	id := utils.AddLock(message)
	link := "https://lock-go.herokuapp.com/" + id
	templates.ExecuteTemplate(w, "new.html", link)
}

func keyGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	data, err := utils.GetLock(key)
	id := data.LockID
	message := data.Message
	if err != nil {
		templates.ExecuteTemplate(w, "404.html", "")
	} else {
		go func() {
			time.Sleep(time.Second * 5)
			utils.DeleteLock(id)

		}()
		templates.ExecuteTemplate(w, "template.html", message)
	}
}

func main() {
	defer db.Close()

	db.AutoMigrate(&utils.Lock{})

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/new", indexPostHandler).Methods("POST")
	r.HandleFunc("/{key}", keyGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	LoadTemplates("templates/*.html")
	http.Handle("/", r)
	log.Println("Serving at localhost:8080...")
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
