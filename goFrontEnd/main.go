package main

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
	ErrorMessage string // Add this field
}

type AuthPageData struct {
	Title      string
	Action     string
	ShowLogout bool
}

type RegistrationPageData struct {
	Title        string
	ErrorMessage string
}

//func registerHandler(w http.ResponseWriter, r *http.Request) {
//	data := RegistrationPageData{
//		Title: "Register for PredictAI",
//	}
//
//	if r.Method == "GET" {
//		tmpl := template.Must(template.ParseFiles("register.gohtml"))
//		tmpl.Execute(w, data)
//	} else if r.Method == "POST" {
//		// Process registration form
//		// If there's an error, set data.ErrorMessage
//		// Re-render the template with the error message
//	}
//}

// DAL needs to be imported and implemented

func connectToDal() {
	// Import dal and connect Dal to the login/register page.
	//http.HandleFunc("/login", loginHandler)
	//http.HandleFunc("/register", registerHandler)
	//http.HandleFunc("/logout", logoutHandler)
	//	http.HandleFunc("/login/google", googleLoginHandler) // Placeholder for Google login
	//	http.HandleFunc("/admin/login", adminLoginHandler)
}

// The functions below are just templates for DAL
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Render the login page
		tmpl := template.Must(template.ParseFiles("path/to/login.gohtml"))
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		// Process the login form
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verify credentials (this is a placeholder, replace with real verification)
		if username == "user" && password == "pass" {
			// Redirect to a secure page or dashboard
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		} else {
			// Render login page with error message
			tmpl := template.Must(template.ParseFiles("path/to/login.gohtml"))
			tmpl.Execute(w, map[string]interface{}{
				"ErrorMessage": "Invalid credentials",
			})
		}
	}
}

//
//func registerHandler(w http.ResponseWriter, r *http.Request) {
//	data := AuthPageData{
//		Title:      "Register for PredictAI",
//		Action:     "register",
//		ShowLogout: false,
//	}
//	tmpl := template.Must(template.ParseFiles("templates/auth.gohtml"))
//	err := tmpl.ExecuteTemplate(w, "register", data)
//	if err != nil {
//		log.Printf("Error executing template: %v", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	}
//}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
	// Redirect or inform the user they've been logged out
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

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
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
