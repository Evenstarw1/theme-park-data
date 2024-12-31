package themepark

type Store interface {
	AddUser(user User) error
	GetUser(userID int) (*User, error)
	GetAllUsers() ([]User, error) // Only for admin users
	UpdateUser(user User) error
	AddThemePark(themePark ThemePark) error    // Only for admin users
	UpdateThemePark(themePark ThemePark) error // Only for admin users
	DeleteThemePark(id int) error              // Only for admin users
	GetAllThemeParks() ([]ThemePark, error)
	GetThemeParkDetail(parkID int) (*ThemePark, error)
	InsertParkComment(parkID int, userId int, comment string) error
	AddCategory(name string) error
	GetAllCategories() ([]Category, error)
	//RemoveCategory() error
	// SignIn returns a token if ok, error if nok
	SignIn(userName string, userPass string) (string, int, error) // Token, userId, error
	// IsLoggedIn given a token, returns if the user is logged in, the user ID and the admin level (1 Admin, 2 normal user).
	IsLoggedIn(token string) (bool, int, int)
}
