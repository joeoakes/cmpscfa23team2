package main

import (
	"cmpscfa23team2/dal"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type AuthData struct {
	Title      string
	Action     string
	ShowLogout bool
	Username   string
	Success    bool
	Error      string
}

type RegistrationPageData struct {
	Title        string
	ErrorMessage string
}

func serveTemplate(content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emptyData := &AuthData{} // Change to AuthData
		renderTemplate(w, "layout.gohtml", emptyData)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *AuthData) {
	t, err := template.ParseFiles("path/to/" + tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if data == nil {
		data = &AuthData{}
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := AuthData{Action: "login"}
		renderTemplate(w, "path/to/login.gohtml", &data)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		log.Printf("Request Headers: %v", r.Header)

		// Call the DAL authentication function
		token, err := dal.AuthenticateUser(username, password)
		if err != nil {
			// Log the authentication error
			log.Printf("Authentication error: %v", err)

			// Handle authentication error
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid credentials"}`))
			return
		}

		// Log the token
		log.Printf("Token: %v", token)

		// Respond with a success message in JSON format
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Authentication successful"}`))

		log.Printf("Response Headers: %v", w.Header())
		return
	}
}

func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin := true
		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	err = dal.LogoutUser(userID)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logout successful"))
}

func extractUserIDFromToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("Authorization header not found")
	}

	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("Invalid Authorization header format")
	}

	tokenString := strings.TrimSpace(splitToken[1])

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(dal.SECRET_KEY), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Failed to parse token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("User ID not found in token claims")
	}

	return userID, nil
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

// Updated requireAdmin middleware to use DAL function
//func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Extract token from the Authorization header
//		header := r.Header.Get("Authorization")
//		if header == "" {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		// Extract token from "Bearer <token>"
//		splitToken := strings.Split(header, "Bearer ")
//		if len(splitToken) != 2 {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		tokenString := strings.TrimSpace(splitToken[1])
//
//		// Validate the token using the DAL function
//		valid, err := dal.ValidateToken(tokenString)
//		if err != nil || !valid {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	}
//}

// Old method to test the dashboard.
// Authentication Middleware to check if the user is logged in and has admin role
//func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Placeholder for actual authentication and role check
//		// You should replace this with a call to your DAL methods to check for a valid admin session/token
//		// For example: isAdmin, err := dal.IsUserAdmin(session.UserID)
//		isAdmin := true // For demonstration purposes, assign false to see the difference
//
//		if !isAdmin {
//			http.Error(w, "Forbidden", http.StatusForbidden)
//			return
//		}
//		next.ServeHTTP(w, r)
//	}
//}

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
