package utils

// Role constants
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
)

// Image upload constants
const (
	ImageUploadPath = "./uploads/images/"
	MaxImageSize    = 10 * 1024 * 1024 // 10MB
)

// Context keys
const (
	ContextUserID   = "userID"
	ContextUserRole = "userRole"
)

// Response status messages
const (
	StatusSuccess = "success"
	StatusError   = "error"
)
