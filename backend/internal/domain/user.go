package domain

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	BaseModel
	Email        string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	FullName     string `gorm:"size:255" json:"full_name"`
	Role         Role   `gorm:"size:32;not null;default:user;index" json:"role"`
	IsActive     bool   `gorm:"not null;default:true" json:"is_active"`
}

func (User) TableName() string { return "users" }
