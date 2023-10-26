package DAL

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// GetUserRole fetches the role associated with a given user ID.
func GetUserRole(userID string) (string, error) {
	var userRole string
	err := DB.QueryRow("Call get_user_role(?)", userID).Scan(&userRole)
	if err != nil {
		log.Printf("Error in GetUserRole: %v", err)
		return "", err
	}
	log.Printf("GetUserRole: User Role for UserID %s is %s", userID, userRole)
	return userRole, nil
}

// IsUserActive checks if a user is currently marked as active based on their user ID.
func IsUserActive(userID string) (bool, error) {
	var isActive bool
	err := DB.QueryRow("CALL is_user_active(?)", userID).Scan(&isActive)
	if err != nil {
		log.Printf("Error in IsUserActive: %v", err)
		return false, err
	}
	log.Printf("IsUserActive: UserID %s is active: %v", userID, isActive)
	return isActive, nil
}

// AuthorizeUser verifies if a user has the necessary role to perform a certain action.
func AuthorizeUser(userID string, requiredRole string) (bool, error) {
	userRole, err := GetUserRole(userID)
	if err != nil {
		log.Printf("Error in AuthorizeUser: %v", err)
		return false, err
	}
	hasPermission := userRole == requiredRole
	log.Printf("AuthorizeUser: UserID %s has required role %s: %v", userID, requiredRole, hasPermission)
	return hasPermission, nil
}

// GetPermissionsForRole fetches all permissions associated with a given user role.
func GetPermissionsForRole(userRole string) ([]Permission, error) {
	// Execute a stored procedure to fetch permissions for the user role.
	rows, err := DB.Query("CALL get_permissions_for_role(?)", userRole)
	if err != nil {
		log.Printf("Error in GetPermissionsForRole: %v", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var action, resource string
		if err := rows.Scan(&action, &resource); err != nil {
			log.Printf("Error in GetPermissionsForRole (Scan): %v", err)
			return nil, err
		}
		permissions = append(permissions, NewPermission(action, resource))
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error in GetPermissionsForRole (Rows): %v", err)
		return nil, err
	}

	log.Printf("GetPermissionsForRole: Permissions for Role %s: %+v", userRole, permissions)
	return permissions, nil
}

// CheckPermission verifies if a specific role has permission to perform a certain action on a given resource.
func CheckPermission(userRole, action, resource string) (bool, error) {
	// Execute a stored procedure to check if the role has the permission.
	var hasPermission bool
	err := DB.QueryRow("CALL check_permission(?, ?, ?)", userRole, action, resource).Scan(&hasPermission)
	if err != nil {
		log.Printf("Error in CheckPermission: %v", err)
		return false, err
	}
	log.Printf("CheckPermission: Role %s has permission for Action %s on Resource %s: %v", userRole, action, resource, hasPermission)
	return hasPermission, nil
}

// UpdateUserRole allows for changing the role associated with a user.
func UpdateUserRole(userID, newRole string) error {
	_, err := DB.Exec("CALL update_user_role(?, ?)", userID, newRole)
	if err != nil {
		log.Printf("Error in UpdateUserRole: %v", err)
	} else {
		log.Printf("UpdateUserRole: Role updated for UserID %s to %s", userID, newRole)
	}
	return err
}

// DeactivateUser marks a user as inactive.
func DeactivateUser(userID string) error {
	_, err := DB.Exec("CALL deactivate_user(?)", userID)
	if err != nil {
		log.Printf("Error in DeactivateUser: %v", err)
	} else {
		log.Printf("DeactivateUser: UserID %s marked as inactive", userID)
	}
	return err
}

// AddPermission allows for adding a new permission to a user role.
func AddPermission(userRole, action, resource string) error {
	_, err := DB.Exec("CALL add_permission(?, ?, ?)", userRole, action, resource)
	if err != nil {
		log.Printf("Error in AddPermission: %v", err)
	} else {
		log.Printf("AddPermission: Permission added for Role %s: Action %s on Resource %s", userRole, action, resource)
	}
	return err
}

// HasPermission is a higher-level function to check if a user has a specific permission.
func HasPermission(userID, action, resource string) (bool, error) {
	userRole, err := GetUserRole(userID)
	if err != nil {
		log.Printf("Error in HasPermission (GetUserRole): %v", err)
		return false, err
	}

	hasPermission, err := CheckPermission(userRole, action, resource)
	if err != nil {
		log.Printf("Error in HasPermission (CheckPermission): %v", err)
		return false, err
	}

	log.Printf("HasPermission: UserID %s has permission for Action %s on Resource %s: %v", userID, action, resource, hasPermission)
	return hasPermission, nil
}

// ... [Additional functions for other sprocs as required]

// Permission represents a user's permission to perform an action on a resource.
type Permission struct {
	Action   string
	Resource string
}

// NewPermission creates a new Permission object.
func NewPermission(action, resource string) Permission {
	return Permission{Action: action, Resource: resource}
}

// ...
