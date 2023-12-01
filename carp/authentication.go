package dal

//
//import (
//	"cmpscfa23team2/dal"
//	"encoding/json"
//	"log"
//	"net/http"
//)
//
//// LoginHandler handles the user login process.
//func LoginHandler(w http.ResponseWriter, r *http.Request) {
//	// Retrieve username and password from the request
//	username := r.FormValue("username")
//	password := r.FormValue("password")
//
//	// Call the DAL authentication function
//	token, err := dal.AuthenticateUser(username, password)
//	if err != nil {
//		// Log the authentication error
//		log.Printf("Authentication error: %v", err)
//
//		// Handle authentication error
//		w.WriteHeader(http.StatusUnauthorized)
//		w.Write([]byte(`{"error": "Invalid credentials"}`))
//		return
//	}
//
//	// Log the token
//	log.Printf("Token: %v", token)
//
//	// Respond with a success message in JSON format
//	w.Header().Set("Content-Type", "application/json")
//	w.Write([]byte(`{"message": "Authentication successful"}`))
//	return
//}
//
//// writeJSONResponse writes a JSON response with the provided data.
//func writeJSONResponse(w http.ResponseWriter, data interface{}) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//
//	if err := json.NewEncoder(w).Encode(data); err != nil {
//		log.Printf("Error encoding JSON response: %v", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	}
//}
//
//// LogoutHandler handles the user logout process.
//func LogoutHandler(w http.ResponseWriter, r *http.Request) {
//	// Retrieve user ID from the request, depending on your authentication mechanism
//	userID := r.FormValue("userID")
//
//	// Call the DAL logout function to log out the user
//	err := dal.LogoutUser(userID)
//	if err != nil {
//		// Handle logout error
//		http.Error(w, "Logout failed", http.StatusInternalServerError)
//		return
//	}
//
//	w.Write([]byte("Logout successful"))
//}
