package carp

import (
	"cmpscfa23team2/dal"
	"net/http"
)

// LoginHandler handles the user login process.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve username and password from the request
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Call the DAL authentication function to authenticate the user
	token, err := dal.AuthenticateUser(username, password)
	if err != nil {
		// Handle authentication error
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Set the token in the response header or send it as a JSON response, as needed
	w.Header().Set("Authorization", "Bearer "+token)
	w.Write([]byte("Login successful"))
}

// LogoutHandler handles the user logout process.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from the request, depending on your authentication mechanism
	userID := r.FormValue("userID")

	// Call the DAL logout function to log out the user
	err := dal.LogoutUser(userID)
	if err != nil {
		// Handle logout error
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	// Optionally, invalidate the token or perform other logout-related tasks
	// ...

	w.Write([]byte("Logout successful"))
}

// Changes Made:
// 1. Created a package named "carp" to encapsulate login and logout functionality.
// 2. Utilized the DAL package for user authentication and logout operations.
// 3. Added comments to describe the purpose and functionality of each function.
// 4. Implemented error handling for authentication and logout processes.
// 5. Set the JWT token in the response header upon successful user authentication.
// 6. Included optional steps for token invalidation or other logout-related tasks.
// 7. Incorporated consistent naming conventions for functions and variables.
// 8. Ensured that appropriate HTTP status codes are used for error responses.
// 9. Left placeholders for retrieving user ID based on the authentication mechanism.
// 10. Encapsulated the login and logout functionalities within the "carp" package.
