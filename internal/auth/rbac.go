package auth

import "errors"

// Common roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// Permission represents a specific action that can be performed
type Permission string

// Define permissions
const (
	PermissionReadNote      Permission = "read:note"
	PermissionCreateNote    Permission = "create:note"
	PermissionUpdateNote    Permission = "update:note"
	PermissionDeleteNote    Permission = "delete:note"
	PermissionManageUsers   Permission = "manage:users"
	PermissionViewAnalytics Permission = "view:analytics"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidRole  = errors.New("invalid role")
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[string][]Permission{
	RoleUser: {
		PermissionReadNote,
		PermissionCreateNote,
		PermissionUpdateNote,
		PermissionDeleteNote,
	},
	RoleAdmin: {
		PermissionReadNote,
		PermissionCreateNote,
		PermissionUpdateNote,
		PermissionDeleteNote,
		PermissionManageUsers,
		PermissionViewAnalytics,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role string, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// ValidateRole checks if a role is valid
func ValidateRole(role string) error {
	if role != RoleUser && role != RoleAdmin {
		return ErrInvalidRole
	}
	return nil
}
