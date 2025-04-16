package model

type UserRole string

const (
	EmployeeRole  UserRole = "employee"
	ModeratorRole UserRole = "moderator"
)

type User struct {
	ID    string `json:"id,omitempty" format:"uuid"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=employee moderator"`
}
