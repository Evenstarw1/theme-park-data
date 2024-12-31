package themepark

import "time"

type User struct {
	ID             int       `json:"id,omitempty"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"`
	AccessLevel    int       `json:"access_level,omitempty"`
	BirthDate      time.Time `json:"birth_date"`
	City           string    `json:"city"`
	Categories     []int     `json:"categories"`
	ProfilePicture string    `json:"profile_picture"`
	Description    string    `json:"description"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	UserId      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type Category struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type UsersCategory struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	CategoryId int       `json:"category_id"`
	Created    time.Time `json:"created"`
}

type ThemeparkCategory struct {
	Id          int       `json:"id"`
	ThemeparkId int       `json:"themepark_id"`
	CategoryId  int       `json:"category_id"`
	Created     time.Time `json:"created"`
}

type Comment struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	UserName    string    `json:"user_name,omitempty"`
	ThemeparkId int       `json:"themepark_id,omitempty"`
	Comment     string    `json:"comment"`
	Created     time.Time `json:"created"`
}

type Attraction struct {
	Id          int       `json:"id"`
	ThemeparkId int       `json:"themepark_id"`
	Name        string    `json:"name"`
	Created     time.Time `json:"created"`
}

type ThemePark struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Picture     string       `json:"picture"`
	Latitude    float64      `json:"latitude,omitempty"`
	Longitude   float64      `json:"longitude,omitempty"`
	Categories  []Category   `json:"categories,omitempty"`
	Comments    []Comment    `json:"comments,omitempty"`
	Attractions []Attraction `json:"attractions,omitempty"`
}
