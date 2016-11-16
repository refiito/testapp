package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type Answer struct {
	Email   string
	Message string
}

var emailRegex = regexp.MustCompile(`^[^<>\\#$@\s]+@[^<>\\#$@\s]*[^<>\\#$\.\s@]{1}?\.{1}?[^<>\\#$\.@\s]{1}?[^<>\\#$@\s]+$`)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	u := &Answer{Message: "", Email: "E-mail address"}
	renderTemplate(w, "index", u)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	message := fmt.Sprintf("Hi, %s Thanks for signing up, now log in using your e-mail and password", email)
	if !emailRegex.Match([]byte(email)) {
		message = fmt.Sprintf("Invalid e-mail: '%s'", email)
	}

	err := createUser(email, password)
	if err != nil {
		message = err.Error()
	}

	u := &Answer{
		Email:   email,
		Message: message,
	}
	renderTemplate(w, "index", u)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	u := &Answer{
		Email:   email,
		Message: fmt.Sprintf("Hi, %s, Hello Again!", email),
	}

	user, err := getUser(email)
	if err != nil {
		u.Message = err.Error()
		renderTemplate(w, "index_result", u)
		return
	}

	err2 := user.authenticate(password)
	if err != nil {
		u.Message = err2.Error()
	}
	renderTemplate(w, "index_result", u)
}

var templates = template.Must(template.ParseFiles("index.html", "index_result.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, u *Answer) {
	err := templates.ExecuteTemplate(w, tmpl+".html", u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	configure()
	runDBCheck()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	server := &http.Server{
		Addr: getListenAddress(),
	}

	log.Printf("server starting at %s\n", getListenAddress())
	log.Fatal(server.ListenAndServe())
}
