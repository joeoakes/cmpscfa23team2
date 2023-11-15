package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// DAL needs to be imported and implemented

func connectToDal() {
	// Import dal and connect Dal to the login/register page.
	//http.HandleFunc("/login", loginHandler)
	//http.HandleFunc("/register", registerHandler)
}

func main() {

	//connectToDal()

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "PredictAI",
			Content: "home", // Indicates which content template to use
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "About Us - PredictAI",
			Content: "about",
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/contributors", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "Our Contributors - PredictAI",
			Content: "contributors",
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "Login - PredictAI",
			Content: "login",
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Serving static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
