package models

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username" unique:"true"` // Assuming database supports unique constraints
	Email    string `json:"email" unique:"true"`    // Assuming database supports unique constraints
	Password string `json:"password"`               // Include password for storage
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}
