package domain

type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
