package user

import "time"

type User struct {
	ID              int64     `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           string    `json:"phone,omitempty"`
	DOB             string    `json:"dob,omitempty"`
	Address         string    `json:"address,omitempty"`
	Email           string    `json:"email"`
	Password        string    `json:"password,omitempty"`
	ConfirmPassword string    `json:"confirm_password,omitempty"`
	PasswordHash    string    `json:"-"`
	Role            string    `json:"role,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	AvatarURL 		string 	  `json:"avatar_url"`
}