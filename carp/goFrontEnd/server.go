package main

// The errors: the file does not see the other methods from the middleware
// change the configurations of the build to: from 'file' to 'directory' goFrontEnd
// This is how the filepath looks for me:
// C:\Users\Public\GoLandProjects\PredictAi\carp\goFrontEnd

import (
	"cmpscfa23team2/dal"
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)
	setupServer()
}

func setupServer() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)

	files, err := filepath.Glob(filepath.Join(dir, "templates/*.gohtml"))
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
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler(tmpl, w, r)
	})
	//http.HandleFunc("/register", makeHandler(tmpl, "register"))
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

func registerHandler(tmpl *template.Template, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	switch r.Method {
	case "GET":
		// Display the registration form
		err := tmpl.ExecuteTemplate(w, "register", RegistrationPageData{Title: "register"})
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	case "POST":

		// Parse form values
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Extract form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")

		// Check if passwords match
		if password != confirmPassword {
			tmpl.ExecuteTemplate(w, "register", RegistrationPageData{
				Title:        "Register",
				ErrorMessage: "Passwords do not match",
			})
			return
		}

		// Set default values for role and active status
		defaultRole := "USR" // Modify as necessary
		active := true       // Set to false if you require email verification, etc.

		// Call DAL function to register user
		_, err := dal.RegisterUser(username, email, defaultRole, password, active)
		if err != nil {
			tmpl.ExecuteTemplate(w, "register", RegistrationPageData{
				Title:        "Register",
				ErrorMessage: "Registration failed: " + err.Error(),
			})
			return
		}

		// Redirect on successful registration
		http.Redirect(w, r, "/login", http.StatusFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
