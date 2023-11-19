package main

// Need to complete and connect this file with DAL
// Need to make the index, login, register functional through DAL
// Play with the requireAdmin
// Connect home page with CUDA and DAL

import (
	"html/template"
	"log"
	"net/http"
)

type AuthPageData struct {
	Title      string
	Action     string
	ShowLogout bool
}

type RegistrationPageData struct {
	Title        string
	ErrorMessage string
}

// DAL needs to be imported and implemented

func serveTemplate(content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "layout.gohtml", PageData{Title: "PredictAI - " + content, Content: content})
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles("path/to/" + tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
	// Redirect or inform the user they've been logged out
}

/*
Backend:
Data Retrieval:
Implement functions in your DAL that retrieve the necessary data to populate the dashboard and settings page. (Already done in DAL)
This may include user data, system configurations, etc.

Data Update:
Create functions in your DAL to handle updates to the system configurations. (Already done in DAL)
These will be used when the settings form is submitted.

API Endpoints:
Define API endpoints in your Go server to handle GET requests for loading data into
the dashboard and POST requests for updating settings.
*/

/*
Frontend:

Settings Form:
Ensure the settings form is correctly linked to your backend.
Use name attributes in your form inputs to capture the data in your Go backend.

Dynamic Data Loading:
Use JavaScript or this file go methods to dynamically load data into the dashboard if needed.

Form Submission Handling: Write JavaScript to handle the form submission asynchronously (AJAX) to provide a smoother user experience without reloading the page.
*/

// Example for GET handler to load dashboard data
//http.HandleFunc("/api/dashboard", func(w http.ResponseWriter, r *http.Request) {
//	users := dal.GetUsers() // Assuming you have a function to get user data
//	jsonResponse, _ := json.Marshal(users)
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(jsonResponse)
//})

// Example for POST handler to update settings
//http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		// Parse form values
//		r.ParseForm()
//		crawlingRules := r.FormValue("crawlingRules")
//		dataMapping := r.FormValue("dataMapping")
//
//		// Update settings in DAL
//		dal.UpdateCrawlingRules(crawlingRules)
//		dal.UpdateDataMapping(dataMapping)
//
//		// Redirect or send a success response
//		http.Redirect(w, r, "/settings", http.StatusFound)
//	}
//})

// Authentication Middleware to check if the user is logged in and has admin role
func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Placeholder for actual authentication and role check
		// You should replace this with a call to your DAL methods to check for a valid admin session/token
		// For example: isAdmin, err := dal.IsUserAdmin(session.UserID)
		isAdmin := true // For demonstration purposes, assign false to see the difference

		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

//http.HandleFunc("/start-crawler", func(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	// Call the StartCrawler function from your dal package
//	err := dal.StartCrawler()
//	if err != nil {
//		http.Error(w, "Failed to start crawler", http.StatusInternalServerError)
//		return
//	}
//	// Respond with a success message
//	json.NewEncoder(w).Encode(map[string]string{"message": "Crawler started"})
//})

// Define similar handlers for stopping the crawler, viewing logs, etc.
