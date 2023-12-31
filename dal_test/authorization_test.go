package dal_test

import (
	"cmpscfa23team2/dal"
	"testing"
)

func TestAddPermission(t *testing.T) {
	userRole := "USR"           // Replace with a valid user role
	action := "READ"            // Replace with a valid action
	resource := "SOME_RESOURCE" // Replace with a valid resource
	err := dal.AddPermission(userRole, action, resource)
	if err != nil {
		dal.InsertLog("400", "Failed to add permission", "TestAddPermission()")
		if err != nil {
			return
		}
		t.Fatalf("Failed to add permission: %v", err)
	}

	// Add assertions to verify the permission is added in the database.
	hasPermission, err := dal.CheckPermission(userRole, action, resource)
	if err != nil {
		dal.InsertLog("400", "Failed to check added permission", "TestAddPermission()")
		if err != nil {
			return
		}
		t.Fatalf("Failed to check added permission: %v", err)
	}

	if !hasPermission {
		dal.InsertLog("400", "Permission not found in database", "TestAddPermission()")
		if err != nil {
			return
		}
		t.Errorf("Expected permission to be added, but it was not found in the database")
	} else {
		dal.InsertLog("200", "Successfully added and verified permission", "TestAddPermission()")
		if err != nil {
			return
		}
	}
}
func TestGetUserRole(t *testing.T) {
	userID := "7e8e9aa4-8f2c-11ee-ae02-30d042e80ac3" // Replace with a valid user ID
	role, err := dal.GetUserRole(userID)
	if err != nil {
		dal.InsertLog("400", "Failed to get user role", "TestGetUserRole()")
		if err != nil {
			return
		}
		t.Fatalf("Failed to get user role: %v", err)
	}

	// Assert the role is as expected
	expectedRole := "ADM" // Replace with the expected role
	if role != expectedRole {
		dal.InsertLog("400", "Role mismatch", "TestGetUserRole()")
		t.Errorf("Expected role %s, but got %s", expectedRole, role)
	} else {
		dal.InsertLog("200", "Successfully got user role", "TestGetUserRole()")
	}
}

func TestIsUserActive(t *testing.T) {
	userID := "7e8e9aa4-8f2c-11ee-ae02-30d042e80ac3" // Replace with a valid user ID
	isActive, err := dal.IsUserActive(userID)
	if err != nil {
		dal.InsertLog("400", "Failed to check user's activity status", "TestIsUserActive()")
		t.Fatalf("Failed to check user's activity status: %v", err)
	}

	// Assert the isActive status is as expected
	expectedIsActive := true // Replace with the expected value
	if isActive != expectedIsActive {
		dal.InsertLog("400", "Active status mismatch", "TestIsUserActive()")
		t.Errorf("Expected IsActive %t, but got %t", expectedIsActive, isActive)
	} else {
		dal.InsertLog("200", "Successfully checked user's activity status", "TestIsUserActive()")
	}
}

func TestAuthorizeUser(t *testing.T) {
	userID := "7e8e9aa4-8f2c-11ee-ae02-30d042e80ac3" // Replace with a valid user ID
	requiredRole := "ADM"                            // Replace with the required role
	isAuthorized, err := dal.AuthorizeUser(userID, requiredRole)
	if err != nil {
		dal.InsertLog("400", "Failed to authorize user", "TestAuthorizeUser()")
		t.Fatalf("Failed to authorize user: %v", err)
	}

	// Assert the authorization status is as expected
	expectedAuthorization := true // Replace with the expected value
	if isAuthorized != expectedAuthorization {
		dal.InsertLog("400", "Authorization mismatch", "TestAuthorizeUser()")
		t.Errorf("Expected authorization %t, but got %t", expectedAuthorization, isAuthorized)
	} else {
		dal.InsertLog("200", "Successfully authorized user", "TestAuthorizeUser()")
	}
}

func TestGetPermissionsForRole(t *testing.T) {
	userRole := "ADM" // Replace with a valid user role
	permissions, err := dal.GetPermissionsForRole(userRole)
	if err != nil {
		dal.InsertLog("400", "Failed to get permissions for role", "TestGetPermissionsForRole()")
		t.Fatalf("Failed to get permissions for role: %v", err)
	}

	// Assert that the permissions slice is not empty or nil
	if permissions == nil || len(permissions) == 0 {
		dal.InsertLog("400", "No permissions retrieved", "TestGetPermissionsForRole()")
		t.Errorf("Expected non-empty permissions slice, but got empty")
	} else {
		dal.InsertLog("200", "Successfully retrieved permissions for role", "TestGetPermissionsForRole()")
	}
}

func TestCheckPermission(t *testing.T) {
	userRole := "ADM"           // Replace with a valid user role
	action := "READ"            // Replace with a valid action
	resource := "SOME_RESOURCE" // Replace with a valid resource
	hasPermission, err := dal.CheckPermission(userRole, action, resource)
	if err != nil {
		dal.InsertLog("400", "Failed to check permission", "TestCheckPermission()")
		t.Fatalf("Failed to check permission: %v", err)
	}

	// Assert the permission status is as expected
	expectedPermission := true // Replace with the expected value
	if hasPermission != expectedPermission {
		dal.InsertLog("400", "Permission mismatch", "TestCheckPermission()")
		t.Errorf("Expected permission %t, but got %t", expectedPermission, hasPermission)
	} else {
		dal.InsertLog("200", "Successfully checked permission", "TestCheckPermission()")
	}
}
func TestHasPermission(t *testing.T) {
	userID := "7e8e9aa4-8f2c-11ee-ae02-30d042e80ac3" // Replace with a valid user ID
	action := "READ"                                 // Replace with a valid action
	resource := "SOME_RESOURCE"                      // Replace with a valid resource
	hasPermission, err := dal.HasPermission(userID, action, resource)
	if err != nil {
		dal.InsertLog("400", "Failed to check permission", "TestHasPermission()")
		t.Fatalf("Failed to check permission: %v", err)
	}

	// Assert the permission status is as expected
	expectedPermission := true // Replace with the expected value
	if hasPermission != expectedPermission {
		dal.InsertLog("400", "Permission mismatch", "TestHasPermission()")
		t.Errorf("Expected permission %t, but got %t", expectedPermission, hasPermission)
	} else {
		dal.InsertLog("200", "Successfully checked permission", "TestHasPermission()")
	}
}
func TestUpdateUserRole(t *testing.T) {
	userID := "7e8ec5d9-8f2c-11ee-ae02-30d042e80ac3" // Replace with a valid user ID
	newRole := "DEV"                                 // Replace with the new role
	err := dal.UpdateUserRole(userID, newRole)
	if err != nil {
		dal.InsertLog("400", "Failed to update user role", "TestUpdateUserRole()")
		t.Fatalf("Failed to update user role: %v", err)
	}

	// Add assertions to verify the role has been updated in the database.
	updatedRole, err := dal.GetUserRole(userID)
	if err != nil {
		dal.InsertLog("400", "Failed to get updated user role", "TestUpdateUserRole()")
		t.Fatalf("Failed to get updated user role: %v", err)
	}

	if updatedRole != newRole {
		dal.InsertLog("400", "Updated role mismatch", "TestUpdateUserRole()")
		t.Errorf("Expected updated role %s, but got %s", newRole, updatedRole)
	} else {
		dal.InsertLog("200", "Successfully updated and verified user role", "TestUpdateUserRole()")
	}
}

func TestDeactivateUser(t *testing.T) {
	userID := "7e8ec5d9-8f2c-11ee-ae02-30d042e80ac3" // Replace with a valid user ID
	err := dal.DeactivateUser(userID)
	if err != nil {
		dal.InsertLog("400", "Failed to deactivate user", "TestDeactivateUser()")
		if err != nil {
			return
		}
		t.Fatalf("Failed to deactivate user: %v", err)
	}

	// Add assertions to verify the user is deactivated in the database.
	isActive, err := dal.IsUserActive(userID)
	if err != nil {
		dal.InsertLog("400", "Failed to check user's activity status", "TestDeactivateUser()")
		t.Fatalf("Failed to check user's activity status: %v", err)
	}

	if isActive {
		dal.InsertLog("400", "User is still active after deactivation", "TestDeactivateUser()")
		t.Errorf("Expected user to be deactivated, but user is still active")
	} else {
		dal.InsertLog("200", "Successfully deactivated and verified user", "TestDeactivateUser()")
	}
}
