package dal

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

// Permission represents a user's permission to perform an action on a resource.
//
// The code defines a struct named "Permission" with two fields, "Action" and "Resource," to represent permissions for actions on specific resources.
type Permission struct {
	Action   string
	Resource string
}

// GetUserRole fetches the role associated with a given user ID.
//
// This function retrieves a user's role from a database using the provided userID and logs the result, handling any potential errors.
func GetUserRole(userID string) (string, error) {
	var userRole string
	err := DB.QueryRow("Call get_user_role(?)", userID).Scan(&userRole)
	if err != nil {
		log.Printf("Error in GetUserRole: %v", err)
		InsertLog("400", "Error in GetUserRole: "+err.Error(), "GetUserRole()")
		return "", err
	}
	InsertLog("200", "GetUserRole: User Role for UserID "+userID+" is "+userRole, "GetUserRole()")
	log.Printf("GetUserRole: User Role for UserID %s is %s", userID, userRole)
	return userRole, nil
}

// IsUserActive checks if a user is currently marked as active based on their user ID.
//
// It defines a function "IsUserActive" that checks the activity status of a user in a database and returns a boolean indicating whether the user is active or not, along with an error if any.
func IsUserActive(userID string) (bool, error) {
	var isActive bool
	err := DB.QueryRow("CALL is_user_active(?)", userID).Scan(&isActive)
	if err != nil {
		InsertLog("400", "Error in IsUserActive: "+err.Error(), "IsUserActive()")
		log.Printf("Error in IsUserActive: %v", err)
		return false, err
	}
	InsertLog("200", "IsUserActive: UserID "+userID+" is active: "+strconv.FormatBool(isActive), "IsUserActive()")
	log.Printf("IsUserActive: UserID %s is active: %v", userID, isActive)
	return isActive, nil
}

// AuthorizeUser verifies if a user has the necessary role to perform a certain action.
func AuthorizeUser(userID string, requiredRole string) (bool, error) {
	userRole, err := GetUserRole(userID)
	if err != nil {
		InsertLog("400", "Error in AuthorizeUser: "+err.Error(), "AuthorizeUser()")
		log.Printf("Error in AuthorizeUser: %v", err)
		return false, err
	}
	hasPermission := userRole == requiredRole
	InsertLog("200", "AuthorizeUser: UserID "+userID+" has required role "+requiredRole+": "+strconv.FormatBool(hasPermission), "AuthorizeUser()")
	log.Printf("AuthorizeUser: UserID %s has required role %s: %v", userID, requiredRole, hasPermission)
	return hasPermission, nil
}

// GetPermissionsForRole fetches all permissions associated with a given user role.
//
// This code defines a function that retrieves permissions for a given user role from a database using a stored procedure
// and returns them as a slice of Permission objects while handling potential errors.
func GetPermissionsForRole(userRole string) ([]Permission, error) {
	// Execute a stored procedure to fetch permissions for the user role.
	rows, err := DB.Query("CALL get_permissions_for_role(?)", userRole)
	if err != nil {
		InsertLog("400", "Error in GetPermissionsForRole: "+err.Error(), "GetPermissionsForRole()")
		log.Printf("Error in GetPermissionsForRole: %v", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var action, resource string
		if err := rows.Scan(&action, &resource); err != nil {
			InsertLog("400", "Error in GetPermissionsForRole (Scan): "+err.Error(), "GetPermissionsForRole()")
			log.Printf("Error in GetPermissionsForRole (Scan): %v", err)
			return nil, err
		}
		permissions = append(permissions, NewPermission(action, resource))
	}

	if err := rows.Err(); err != nil {
		InsertLog("400", "Error in GetPermissionsForRole (Rows): "+err.Error(), "GetPermissionsForRole()")
		log.Printf("Error in GetPermissionsForRole (Rows): %v", err)
		return nil, err
	}

	InsertLog("200", "GetPermissionsForRole: Permissions for Role "+userRole+": "+fmt.Sprintf("%+v", permissions), "GetPermissionsForRole()")
	log.Printf("GetPermissionsForRole: Permissions for Role %s: %+v", userRole, permissions)
	return permissions, nil
}

// CheckPermission verifies if a specific role has permission to perform a certain action on a given resource.
func CheckPermission(userRole, action, resource string) (bool, error) {
	// Execute a stored procedure to check if the role has the permission.
	var hasPermission bool
	err := DB.QueryRow("CALL check_permission(?, ?, ?)", userRole, action, resource).Scan(&hasPermission)
	if err != nil {
		InsertLog("400", "Error in CheckPermission: "+err.Error(), "CheckPermission()")
		log.Printf("Error in CheckPermission: %v", err)
		return false, err
	}
	InsertLog("200", "CheckPermission: Role "+userRole+" has permission for Action "+action+" on Resource "+resource+": "+strconv.FormatBool(hasPermission), "CheckPermission()")
	log.Printf("CheckPermission: Role %s has permission for Action %s on Resource %s: %v", userRole, action, resource, hasPermission)
	return hasPermission, nil
}

// UpdateUserRole allows for changing the role associated with a user.
//
// It defines a function UpdateUserRole that updates a user's role in a database using a stored procedure and logs the outcome, handling potential errors.
func UpdateUserRole(userID, newRole string) error {
	_, err := DB.Exec("CALL update_user_role(?, ?)", userID, newRole)
	if err != nil {
		InsertLog("400", "Error in UpdateUserRole: "+err.Error(), "UpdateUserRole()")
		log.Printf("Error in UpdateUserRole: %v", err)
	} else {
		InsertLog("200", "UpdateUserRole: Role updated for UserID "+userID+" to "+newRole, "UpdateUserRole()")
		log.Printf("UpdateUserRole: Role updated for UserID %s to %s", userID, newRole)
	}
	return err
}

// DeactivateUser marks a user as inactive.
//
// It deactivates a user in a database by calling a stored procedure with the provided userID and logs the outcome, handling any errors that may occur.
func DeactivateUser(userID string) error {
	_, err := DB.Exec("CALL deactivate_user(?)", userID)
	if err != nil {
		InsertLog("400", "Error in DeactivateUser: "+err.Error(), "DeactivateUser()")
		log.Printf("Error in DeactivateUser: %v", err)
	} else {
		InsertLog("200", "DeactivateUser: UserID "+userID+" marked as inactive", "DeactivateUser()")
		log.Printf("DeactivateUser: UserID %s marked as inactive", userID)
	}
	return err
}

// AddPermission allows for adding a new permission to a user role.
func AddPermission(userRole, action, resource string) error {
	_, err := DB.Exec("CALL add_permission(?, ?, ?)", userRole, action, resource)
	if err != nil {
		InsertLog("400", "Error in AddPermission: "+err.Error(), "AddPermission()")
		log.Printf("Error in AddPermission: %v", err)
	} else {
		InsertLog("200", "AddPermission: Permission added for Role "+userRole+": Action "+action+" on Resource "+resource, "AddPermission()")
		log.Printf("AddPermission: Permission added for Role %s: Action %s on Resource %s", userRole, action, resource)
	}
	return err
}

// HasPermission is a higher-level function to check if a user has a specific permission.
//
// It defines a function, HasPermission, which checks if a user has a specific permission by first retrieving the user's role, then verifying the permission for a given action and resource, and logging the result along with potential errors.
func HasPermission(userID, action, resource string) (bool, error) {
	userRole, err := GetUserRole(userID)
	if err != nil {
		InsertLog("400", "Error in HasPermission (GetUserRole): "+err.Error(), "HasPermission()")
		log.Printf("Error in HasPermission (GetUserRole): %v", err)
		return false, err
	}

	hasPermission, err := CheckPermission(userRole, action, resource)
	if err != nil {
		InsertLog("400", "Error in HasPermission (CheckPermission): "+err.Error(), "HasPermission()")
		log.Printf("Error in HasPermission (CheckPermission): %v", err)
		return false, err
	}

	InsertLog("200", "HasPermission: UserID "+userID+" has permission for Action "+action+" on Resource "+resource+": "+strconv.FormatBool(hasPermission), "HasPermission()")
	log.Printf("HasPermission: UserID %s has permission for Action %s on Resource %s: %v", userID, action, resource, hasPermission)
	return hasPermission, nil
}

// NewPermission creates a new Permission object.
//
// This code defines a Go function named "NewPermission" that creates
// and returns a new "Permission" struct with provided "action" and "resource" strings.
func NewPermission(action, resource string) Permission {
	return Permission{Action: action, Resource: resource}
}

// ...
