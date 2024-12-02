package themepark

type Store interface {
	AddUser(user User) error
	GetUser(userID int) (*User, error)
	//GetAllUsers() error
	UpdateUser(user User) error
	//
	//AddThemePark() error
	//GetThemePark() error
	//GetAllThemeParks() error
	//
	//AddComment() error
	//GetComments() error
	//
	//AddCategory() error
	GetAllCategories() ([]Category, error)
	//RemoveCategory() error
	// SignIn returns a token if ok, error if nok
	SignIn(userName string, userPass string) (string, int, error) // Token, userId, error
	IsLoggedIn(token string) bool
	//
}
