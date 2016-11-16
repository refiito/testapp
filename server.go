package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type User struct {
	Email   string
	Message string
}

var emailRegex = regexp.MustCompile(`^[^<>\\#$@\s]+@[^<>\\#$@\s]*[^<>\\#$\.\s@]{1}?\.{1}?[^<>\\#$\.@\s]{1}?[^<>\\#$@\s]+$`)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	u := &User{Message: "", Email: "E-mail address"}
	renderTemplate(w, "index", u)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	message := fmt.Sprintf("Hi, %s Thanks for signing up, now log in using your e-mail and password", email)
	if !emailRegex.Match([]byte(email)) {
		message = fmt.Sprintf("Invalid e-mail: '%s'", email)
	}

	u := &User{
		Email:   email,
		Message: message,
	}
	renderTemplate(w, "index", u)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// email := r.FormValue("email")
	// password := r.FormValue("password")
	u := &User{Email: "foo"}
	renderTemplate(w, "index_result", u)
}

var templates = template.Must(template.ParseFiles("index.html", "index_result.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, u *User) {
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
