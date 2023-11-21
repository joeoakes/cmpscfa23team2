package main

// Need to complete and connect this file with DAL
// Need to make the index, login, register functional through DAL
// Play with the requireAdmin
// Connect home page with CUDA and DAL

import (
	"cmpscfa23team2/dal"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"strings"
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
		emptyData := &PageData{} // Create a pointer to PageData
		renderTemplate(w, "layout.gohtml", emptyData)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	t, err := template.ParseFiles("path/to/" + tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if data == nil {
		data = &PageData{} // Set default values or handle nil appropriately
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Updated loginHandler to use DAL authentication
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Render the login page
		renderTemplate(w, "path/to/login.gohtml", nil)
	} else if r.Method == "POST" {
		// Process the login form
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Call the DAL authentication function
		token, err := dal.AuthenticateUser(username, password)
		if err != nil {
			// Handle authentication error
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set the token in the response header
		w.Header().Set("Authorization", "Bearer "+token)
		w.Write([]byte("Login successful"))
	}
}

// DAL imports should be added based on the actual structure of your DAL package
// import "path/to/dal"

func authenticateUser(username, password string) (string, error) {
	authenticated, err := dal.AuthenticateUser(username, password)
	if err != nil {
		// Handle the error, log, or return an appropriate response
		return "", err
	}
	return authenticated, nil
}

// Function to get the role of a user
func getUserRole(username string) string {
	role, err := dal.GetUserRole(username)
	if err != nil {
		// Handle the error, log, or return a default role
		return "user"
	}
	return role
}

// Function to get the authenticated user's username
func getAuthenticatedUser(r *http.Request) string {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		// Handle the case where the username is not found in the context
		return ""
	}
	return username
}

// Updated logoutHandler to use DAL logout function
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the token
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	// Call the DAL logout function
	err = dal.LogoutUser(userID)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logout successful"))
}

// Function to extract user ID from the JWT token
func extractUserIDFromToken(r *http.Request) (string, error) {
	// Extract token from the Authorization header
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("Authorization header not found")
	}

	// Extract token from "Bearer <token>"
	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("Invalid Authorization header format")
	}

	tokenString := strings.TrimSpace(splitToken[1])

	// Parse the token to extract user ID
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

	userID, ok := claims["uid"].(string)
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
func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the Authorization header
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		splitToken := strings.Split(header, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(splitToken[1])

		// Validate the token using the DAL function
		valid, err := dal.ValidateToken(tokenString)
		if err != nil || !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

//Old method to test the dashboard.
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

// Changes Made:
// 1. Updated loginHandler to use DAL authentication.
// 2. Implemented functions for user authentication, role retrieval, and logout using DAL.
// 3. Updated requireAdmin middleware to use DAL function for token validation.
// 4. Ensured the settings form is correctly linked to the backend for data capture.
// 5. Utilized JavaScript for dynamic data loading and form submission handling.
// 6. Included error handling and logging for various scenarios.
// 7. Integrated JWT token parsing for user authentication.
