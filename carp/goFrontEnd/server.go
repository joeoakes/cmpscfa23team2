package main

// The errors:
// The file does not 'see' the other methods from the middleware
// change the configurations of the build to: from 'file' to 'directory' goFrontEnd
// This is how the filepath looks for me:
// C:\Users\Public\GoLandProjects\PredictAi\carp\goFrontEnd

import (
	"cmpscfa23team2/dal" // Import the data access layer package
	"encoding/json"      // Import the json package for JSON encoding and decoding
	"html/template"      // Import the template package for HTML templating
	"log"                // Import the log package for logging
	"net/http"           // Import the net/http package for HTTP server and client
	"os"                 // Import the os package for operating system functionality
	"path/filepath"      // Import the filepath package for file path manipulation
)

// PageData struct represents the structure for page data used in templates.
type PageData struct {
	Title        string      // Title of the page
	Content      string      // Content identifier for the page
	ErrorMessage string      // Error message to be displayed on the page
	Users        []*dal.User // Slice of User pointers from the data access layer
}

// main is the entry point of the program.
func main() {
	// Retrieve the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)

	// Set up and start the server
	setupServer()
}

// setupServer configures and starts the HTTP server.
func setupServer() {
	// Retrieve the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current directory:", dir)

	// Glob for template files in the 'templates' directory
	files, err := filepath.Glob(filepath.Join(dir, "templates/*.gohtml"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Template files found:", files)

	// Parse and load the templates
	tmpl := template.Must(template.ParseGlob("templates/*.gohtml"))
	log.Println("Templates loaded:", tmpl.DefinedTemplates())

	// Set up HTTP routes
	setupRoutes(tmpl)

	// Start the server on port 8080
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// setupRoutes configures the HTTP routes for the server.
// tmpl: the parsed templates used for rendering web pages
func setupRoutes(tmpl *template.Template) {
	// Define handlers for various endpoints
	http.HandleFunc("/", makeHandler(tmpl, "login"))
	http.HandleFunc("/home", makeHandler(tmpl, "home"))
	http.HandleFunc("/about", makeHandler(tmpl, "about"))
	http.HandleFunc("/contributors", makeHandler(tmpl, "contributors"))
	http.HandleFunc("/register", makeHandler(tmpl, "register"))
	http.HandleFunc("/documentation", makeHandler(tmpl, "documentation"))
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		dashHandler(tmpl, w, r) // Custom handler for dashboard
	})
	http.HandleFunc("/settings", requireAdmin(makeHandler(tmpl, "settings")))
	http.HandleFunc("/api/predictions", predictionHandler)

	// Serve static files from the 'static' directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

// makeHandler creates an HTTP handler function for a specific content type.
// tmpl: the parsed templates
// content: the content identifier for the page
func makeHandler(tmpl *template.Template, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Special handling for POST requests on register and login
		if r.Method == "POST" && content == "register" {
			registerHandler(tmpl, w, r)
			return
		}
		if r.Method == "POST" && content == "login" {
			loginHandler(tmpl, w, r)
			return
		}

		// Prepare data for the template
		data := struct {
			Title   string
			Content string
		}{
			Title:   "PredictAI - " + content,
			Content: content,
		}

		// Execute the template
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// loginHandler handles login requests and renders the login page.
// tmpl: the parsed templates
// w: the response writer
// r: the HTTP request
func loginHandler(tmpl *template.Template, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Handle GET and POST requests separately
	switch r.Method {
	case "GET":
		renderLoginTemplate(tmpl, w, "")

	case "POST":
		// Process the login form submission
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Authenticate the user
		token, err := dal.AuthenticateUser(email, password)
		if err != nil {
			renderLoginTemplate(tmpl, w, "Invalid email or password")
			return
		}

		// Set the authentication token in a cookie and redirect
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// renderLoginTemplate renders the login page template with an optional error message.
// tmpl: the parsed templates
// w: the response writer
// errorMessage: error message to display on the login page
func renderLoginTemplate(tmpl *template.Template, w http.ResponseWriter, errorMessage string) {
	// Prepare data for the template
	data := PageData{
		Title:        "Login",
		ErrorMessage: errorMessage,
	}

	// Execute the template
	err := tmpl.ExecuteTemplate(w, "login", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// registerHandler handles user registration requests and renders the registration page.
// tmpl: the parsed templates
// w: the response writer
// r: the HTTP request
func registerHandler(tmpl *template.Template, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Handle GET and POST requests separately
	switch r.Method {
	case "GET":
		data := RegistrationPageData{Title: "Register"}
		err := tmpl.ExecuteTemplate(w, "register", data)
		if err != nil {
			log.Printf("Error executing register template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	case "POST":
		// Process the registration form submission
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Extract and validate form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		if password != confirmPassword {
			tmpl.ExecuteTemplate(w, "register", RegistrationPageData{
				Title:        "Register",
				ErrorMessage: "Passwords do not match",
			})
			return
		}

		// Register the user and handle any errors
		_, err := dal.RegisterUser(username, email, "USR", password, true)
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

// predictionHandler handles prediction API requests.
// w: the response writer
// r: the HTTP request
func predictionHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict to GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract and validate query parameters
	domain := r.URL.Query().Get("domain")
	query_identifier := r.URL.Query().Get("queryType")
	if domain == "" || query_identifier == "" {
		http.Error(w, "Missing domain or queryType parameter", http.StatusBadRequest)
		return
	}

	// Fetch and respond with prediction data
	predictionData, err := dal.FetchPredictionData(query_identifier, domain)
	if err != nil {
		log.Printf("Error fetching prediction data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictionData)
}

// renderDashboardTemplate renders the dashboard page template.
// tmpl: the parsed templates
// w: the response writer
// users: slice of User pointers
// errorMessage: error message to display on the dashboard
func renderDashboardTemplate(tmpl *template.Template, w http.ResponseWriter, users []*dal.User, errorMessage string) {
	data := PageData{
		Title:        "Dashboard",
		Users:        users,
		ErrorMessage: errorMessage,
	}
	err := tmpl.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		log.Printf("Error executing dashboard template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// dashHandler handles requests for the dashboard page.
// tmpl: the parsed templates used for rendering the page
// w: the HTTP response writer
// r: the HTTP request received
func dashHandler(tmpl *template.Template, w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to HTML
	w.Header().Set("Content-Type", "text/html")

	// Log the beginning of the handler process
	log.Printf("beginning of dashHandler\n")

	// Retrieve all users using the data access layer
	users, err := dal.GetAllUsers()
	if err != nil {
		// Log and handle any error while fetching users
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Unable to fetch user data", http.StatusInternalServerError)
		return
	}
	log.Printf("before being passed to template: %+v\n")

	// Create a PageData struct with the title, users, and content identifier
	data := PageData{
		Title:   "Dashboard",
		Users:   users,
		Content: "dashboard", // Identifier for the template to render within the layout
	}
	log.Printf("PageData being passed to template: %+v\n", data)

	// Execute the layout template (not the dashboard template directly) with the PageData
	err = tmpl.ExecuteTemplate(w, "layout.gohtml", data)
	if err != nil {
		// Log and handle any error during template execution
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
