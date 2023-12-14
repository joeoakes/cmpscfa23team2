package main

import (
	"cmpscfa23team2/dal"          // Import the data access layer package
	"errors"                      // Import the errors package for error handling
	"github.com/dgrijalva/jwt-go" // Import the jwt-go package for JWT authentication
	"html/template"               // Import the template package for HTML templating
	"log"                         // Import the log package for logging
	"net/http"                    // Import the net/http package for HTTP server and client
	"strings"                     // Import the strings package for string manipulation
)

// AuthData struct represents the authentication data structure.
type AuthData struct {
	Title      string // Title of the page
	Action     string // Action to be performed (like 'login')
	ShowLogout bool   // Flag to show logout option
	Username   string // Username of the logged-in user
	Success    bool   // Flag to indicate successful operation
	Error      string // Error message if any
	LoggedIn   bool   // Added field to track user login status
}

// RegistrationPageData struct represents the registration page data structure.
type RegistrationPageData struct {
	Title        string // Title of the page
	ErrorMessage string // Error message to display on the page
}

// serveTemplate returns an HTTP handler function that renders a template.
// content: the content to be displayed
// loggedIn: a boolean indicating if the user is logged in or not
func serveTemplate(content string, loggedIn bool) http.HandlerFunc {
	emptyData := &AuthData{
		LoggedIn: loggedIn,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "layout.gohtml", emptyData)
	}
}

// renderTemplate renders the specified HTML template.
// w: the response writer
// tmpl: the path to the template file
// data: data to be passed to the template
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

// handleLogin processes login requests.
// w: the response writer
// r: the HTTP request
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Render the login template for GET requests
		data := AuthData{Action: "login"}
		renderTemplate(w, "path/to/login.gohtml", &data)
	} else if r.Method == "POST" {
		// Process POST request for login
		r.ParseForm()
		username := r.FormValue("email") // Update to match the form input name
		password := r.FormValue("password")

		// Call the DAL authentication function
		token, err := dal.AuthenticateUser(username, password)
		if err != nil {
			// Log the authentication error
			log.Printf("Authentication error: %v", err)

			// Render the login template with an error message
			data := AuthData{
				Title:    "Login",
				Action:   "login",
				Error:    "Invalid credentials",
				LoggedIn: false,
			}
			renderTemplate(w, "path/to/login.gohtml", &data)
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

// requireAdmin middleware ensures that the user is an admin before proceeding.
// next: the next handler to call if the user is an admin
func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Placeholder for admin check
		isAdmin := true
		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// logoutHandler handles user logout requests.
// w: the response writer
// r: the HTTP request
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the token
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	// Call the DAL function to log out the user
	err = dal.LogoutUser(userID)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	// Send a successful logout response
	w.Write([]byte("Logout successful"))
}

// extractUserIDFromToken extracts the user ID from the JWT token in the request.
// r: the HTTP request
// Returns the user ID and any error encountered
func extractUserIDFromToken(r *http.Request) (string, error) {
	// Extract the Authorization header
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("Authorization header not found")
	}

	// Split the token from the header
	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("Invalid Authorization header format")
	}

	// Parse the JWT token
	tokenString := strings.TrimSpace(splitToken[1])
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(dal.SECRET_KEY), nil
	})
	if err != nil {
		return "", err
	}

	// Extract the user ID from the token claims
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
