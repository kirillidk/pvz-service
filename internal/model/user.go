package model

type UserRole string

const (
	EmployeeRole  UserRole = "employee"
	ModeratorRole UserRole = "moderator"
)

type User struct {
	ID    string   `json:"id,omitempty" format:"uuid"`
	Email string   `json:"email" binding:"required,email"`
	Role  UserRole `json:"role" binding:"required,oneof=employee moderator"`
}
