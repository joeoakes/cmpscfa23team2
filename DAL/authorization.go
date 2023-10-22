package DAL

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// GetUserRole fetches the role associated with a given user ID.
func GetUserRole(userID string) (string, error) {
	var userRole string
	err := DB.QueryRow("Call get_user_role(?)", userID).Scan(&userRole)
	if err != nil {
		return "", err
	}
	return userRole, nil
}

// IsUserActive checks if a user is currently marked as active based on their user ID.
func IsUserActive(userID string) (bool, error) {
	var isActive bool
	err := DB.QueryRow("CALL is_user_active(?)", userID).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

// AuthorizeUser verifies if a user has the necessary role to perform a certain action.
func AuthorizeUser(userID string, requiredRole string) (bool, error) {
	userRole, err := GetUserRole(userID)
	if err != nil {
		return false, err
	}
	return userRole == requiredRole, nil
}

// GetPermissionsForRole fetches all permissions associated with a given user role.
func GetPermissionsForRole(userRole string) ([]Permission, error) {
	// Execute a stored procedure to fetch permissions for the user role.
	rows, err := DB.Query("CALL get_permissions_for_role(?)", userRole)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var action, resource string
		if err := rows.Scan(&action, &resource); err != nil {
			return nil, err
		}
		permissions = append(permissions, NewPermission(action, resource))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// CheckPermission verifies if a specific role has permission to perform a certain action on a given resource.
func CheckPermission(userRole, action, resource string) (bool, error) {
	// Execute a stored procedure to check if the role has the permission.
	var hasPermission bool
	err := DB.QueryRow("CALL check_permission(?, ?, ?)", userRole, action, resource).Scan(&hasPermission)
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}

// UpdateUserRole allows for changing the role associated with a user.
func UpdateUserRole(userID, newRole string) error {
	_, err := DB.Exec("CALL update_user_role(?, ?)", userID, newRole)
	return err
}

// DeactivateUser marks a user as inactive.
func DeactivateUser(db *sql.DB, userID string) error {
	_, err := DB.Exec("CALL deactivate_user(?)", userID)
	return err
}

// AddPermission allows for adding a new permission to a user role.
func AddPermission(db *sql.DB, userRole, action, resource string) error {
	_, err := DB.Exec("CALL add_permission(?, ?, ?)", userRole, action, resource)
	return err
}

// HasPermission is a higher-level function to check if a user has a specific permission.
func HasPermission(userID, action, resource string) (bool, error) {
	userRole, err := GetUserRole(userID)
	if err != nil {
		return false, err
	}

	return CheckPermission(userRole, action, resource)
}

// Permission represents a user's permission to perform an action on a resource.
type Permission struct {
	Action   string
	Resource string
}

// NewPermission creates a new Permission object.
func NewPermission(action, resource string) Permission {
	return Permission{Action: action, Resource: resource}
}

// ... [Additional functions for other sprocs as required]