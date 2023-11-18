package main

// The errors: the file does not see the other methods from the middleware
// change the configurations of the build to: from 'file' to 'directory' goFrontEnd

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PageData struct {
	Title        string
	Content      string
	ErrorMessage string
}

func main() {
	setupServer()
}

func setupServer() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)

	files, err := filepath.Glob("templates/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Template files found:", files)

	tmpl := template.Must(template.ParseGlob("templates/*.gohtml"))

	setupRoutes(tmpl)

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setupRoutes(tmpl *template.Template) {
	http.HandleFunc("/", makeHandler(tmpl, "home"))
	http.HandleFunc("/about", makeHandler(tmpl, "about"))
	http.HandleFunc("/contributors", makeHandler(tmpl, "contributors"))
	http.HandleFunc("/login", makeHandler(tmpl, "login"))
	http.HandleFunc("/register", registerHandler(tmpl))
	http.HandleFunc("/documentation", makeHandler(tmpl, "documentation"))
	http.HandleFunc("/dashboard", requireAdmin(makeHandler(tmpl, "dashboard")))
	http.HandleFunc("/settings", requireAdmin(makeHandler(tmpl, "settings")))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func makeHandler(tmpl *template.Template, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "PredictAI - " + content,
			Content: content,
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func registerHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Register",
			Content: "register",
			// ErrorMessage can be set based on the context or left empty
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
