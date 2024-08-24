package entity

type Category struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Name     string `json:"name" db:"name"`
}
