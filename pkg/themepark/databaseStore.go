package themepark

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"os"
)

type DatabaseStore struct {
	conn *pgx.Conn
}

func NewDatabaseStore(connectString string) (*DatabaseStore, error) {
	conn, err := pgx.Connect(context.Background(), connectString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	return &DatabaseStore{conn: conn}, nil
}

func (d *DatabaseStore) Close() {
	d.conn.Close(context.Background())
}

func (d *DatabaseStore) GetUser(userID int) (*User, error) {
	var user User
	err := d.conn.QueryRow(context.Background(),
		"select id, name, email, password, access_level, birth_date, city, profile_picture, description from users where id = $1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.AccessLevel, &user.BirthDate, &user.City, &user.ProfilePicture, &user.Description)
	if err != nil {
		return nil, err
	}

	// Get categories for user
	var userCategories []int
	rows, _ := d.conn.Query(context.Background(),
		"select * from users_categories where user_id = $1", userID)

	defer rows.Close()
	for rows.Next() {
		userCategory := UsersCategory{}
		err := rows.Scan(&userCategory.Id, &userCategory.UserId, &userCategory.CategoryId, &userCategory.Created)
		if err != nil {
			return nil, err
		}

		userCategories = append(userCategories, userCategory.CategoryId)
	}

	// Add categories
	user.Categories = userCategories

	return &user, nil
}

func (d *DatabaseStore) UpdateUser(user User) error {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Update user
	_, err = tx.Exec(context.Background(),
		"update users set name = $1, birth_date = $2, city = $3, profile_picture = $4, description = $5 where id = $6",
		user.Name, user.BirthDate, user.City, user.ProfilePicture, user.Description, user.ID)

	if err != nil {
		return err
	}

	// Remove categories
	_, err = tx.Exec(context.Background(),
		"delete from users_categories where user_id = $1", user.ID)

	if err != nil {
		return err
	}

	// Update categories
	_, err = tx.Exec(context.Background(),
		"delete from users_categories where user_id = $1", user.ID)

	if err != nil {
		return err
	}

	for _, category := range user.Categories {
		_, err = tx.Exec(context.Background(),
			`insert into users_categories (user_id, category_id) values ($1,$2)`, user.ID, category)
	}

	tx.Commit(context.Background())
	return nil
}

func (d *DatabaseStore) GetAllUsers() ([]User, error) {
	var users []User

	rows, _ := d.conn.Query(context.Background(),
		"select id, name, email, access_level, birth_date, city, profile_picture, description from users")

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.AccessLevel, &user.BirthDate, &user.City, &user.ProfilePicture, &user.Description)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	rows.Close()

	for i, user := range users {
		// Get categories for user
		var userCategories []int
		rows, _ := d.conn.Query(context.Background(),
			"select * from users_categories where user_id = $1", user.ID)

		for rows.Next() {
			userCategory := UsersCategory{}
			err := rows.Scan(&userCategory.Id, &userCategory.UserId, &userCategory.CategoryId, &userCategory.Created)
			if err != nil {
				return nil, err
			}

			userCategories = append(userCategories, userCategory.CategoryId)
		}

		users[i].Categories = userCategories

		rows.Close()
	}

	return users, nil
}

func (d *DatabaseStore) GetUserFromEmail(email string) (*User, error) {
	var user User
	err := d.conn.QueryRow(context.Background(),
		"select id, name, email, password, access_level, birth_date, city, profile_picture, description from users where email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.AccessLevel, &user.BirthDate, &user.City, &user.ProfilePicture, &user.Description)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *DatabaseStore) SignIn(userEmail string, userPass string) (string, int, error) {
	var user User
	err := d.conn.QueryRow(context.Background(),
		"select id, email, password from users where email = $1", userEmail).
		Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return "", 0, err
	}

	// Check if password matches
	digest := sha256.Sum256([]byte(userPass))
	hashedPass := fmt.Sprintf("%x", digest)
	if hashedPass == user.Password {
		token := uuid.New()
		d.conn.Exec(context.Background(),
			"insert into tokens (token, user_id) values ($1,$2)", token.String(), user.ID)
		return token.String(), user.ID, nil
	}

	return "", 0, errors.New(fmt.Sprintf("Invalid user password for user %s", userEmail))
}

func (d *DatabaseStore) IsLoggedIn(token string) (bool, int, int) {
	var userId int
	var accessLevel int
	rows, _ := d.conn.Query(context.Background(),
		"select u.id, u.access_level from tokens t inner join users u on t.user_id = u.id where t.token = $1", token)

	for rows.Next() { // This means a register exist, meaning the user is logged in
		_ = rows.Scan(&userId, &accessLevel)
		rows.Close()
		return true, userId, accessLevel
	}

	return false, -1, -1
}

func (d *DatabaseStore) GetAllCategories() ([]Category, error) {
	var categories []Category

	rows, _ := d.conn.Query(context.Background(),
		"select * from categories")

	defer rows.Close()
	for rows.Next() {
		category := Category{}
		err := rows.Scan(&category.Id, &category.Name, &category.Created)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (d *DatabaseStore) AddCategory(name string) error {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`insert into categories (name) values ($1)`, name)

	if err != nil {
		return err
	}

	tx.Commit(context.Background())

	return nil
}

func (d *DatabaseStore) AddUser(user User) error {
	// TODO: Check if user already in the db

	passDigest := sha256.Sum256([]byte(user.Password))
	hashedPass := fmt.Sprintf("%x", passDigest)

	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`insert into users (name, email, password, access_level, birth_date, city, profile_picture, description) values ($1,$2,$3,$4,$5,$6,$7,$8)`,
		user.Name, user.Email, hashedPass, 1, user.BirthDate, user.City, user.ProfilePicture, user.Description)

	// Get ID from user
	userWithId, err := d.GetUserFromEmail(user.Email)
	if err != nil {
		return err
	}

	// Insert new users_categories
	for _, category := range user.Categories {
		_, err = tx.Exec(context.Background(),
			`insert into users_categories (user_id, category_id) values ($1,$2)`, userWithId.ID, category)
	}

	tx.Commit(context.Background())

	return nil
}

func (d *DatabaseStore) GetAllThemeParks() ([]ThemePark, error) {
	var themeparks []ThemePark

	rows, _ := d.conn.Query(context.Background(),
		"SELECT id, name, picture FROM themeparks")

	defer rows.Close()
	for rows.Next() {
		themepark := ThemePark{}
		err := rows.Scan(&themepark.Id, &themepark.Name, &themepark.Picture)
		if err != nil {
			rows.Close()
			return nil, err
		}

		themeparks = append(themeparks, themepark)
	}

	rows.Close()

	for i, themepark := range themeparks {
		// Get categories for user
		var parkCategories []Category
		rows, _ := d.conn.Query(context.Background(),
			"select c.id, c.name, c.created from themeparks_categories tc inner join categories c on tc.category_id = c.id where tc.themepark_id = $1", themepark.Id)

		for rows.Next() {
			parkCategory := Category{}
			err := rows.Scan(&parkCategory.Id, &parkCategory.Name, &parkCategory.Created)
			if err != nil {
				rows.Close()
				return nil, err
			}

			parkCategories = append(parkCategories, parkCategory)
		}

		themeparks[i].Categories = parkCategories

		rows.Close()
	}

	return themeparks, nil
}

func (d *DatabaseStore) GetThemeParkDetail(parkID int) (*ThemePark, error) {
	var themepark ThemePark

	err := d.conn.QueryRow(context.Background(),
		"SELECT id, name, description, picture FROM themeparks WHERE id = $1", parkID).
		Scan(&themepark.Id, &themepark.Name, &themepark.Description, &themepark.Picture)
	if err != nil {
		return nil, err
	}

	// Get coordinates
	var point pgtype.Point
	err = d.conn.QueryRow(context.Background(),
		"SELECT location FROM themeparks WHERE id = $1", parkID).
		Scan(&point)
	if err != nil {
		return nil, err
	}

	themepark.Latitude = point.P.X
	themepark.Longitude = point.P.Y

	// Get categories for themeparks
	var categories []Category
	rows, _ := d.conn.Query(context.Background(),
		"select c.id, c.name, c.created from themeparks_categories tc  inner join categories c on tc.category_id = c.id where tc.themepark_id = $1", parkID)

	defer rows.Close()
	for rows.Next() {
		category := Category{}
		err := rows.Scan(&category.Id, &category.Name, &category.Created)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	themepark.Categories = categories

	// Get comments
	var comments []Comment
	rows, _ = d.conn.Query(context.Background(),
		"select c.id, u.id, u.name, c.comment, c.created from comments c inner join users u on c.user_id = u.id where c.themepark_id = $1 order by c.created", parkID)

	defer rows.Close()
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.UserName, &comment.Comment, &comment.Created)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	themepark.Comments = comments

	// Get attractions
	var attractions []Attraction
	rows, _ = d.conn.Query(context.Background(),
		"select * from attractions where themepark_id = $1", parkID)

	defer rows.Close()
	for rows.Next() {
		attraction := Attraction{}
		err := rows.Scan(&attraction.Id, &attraction.ThemeparkId, &attraction.Name, &attraction.Created)
		if err != nil {
			return nil, err
		}

		attractions = append(attractions, attraction)
	}

	themepark.Attractions = attractions

	return &themepark, nil
}
func (d *DatabaseStore) InsertParkComment(parkID int, userId int, comment string) error {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`insert into comments (themepark_id, user_id, comment) values ($1,$2,$3)`,
		parkID, userId, comment)

	tx.Commit(context.Background())

	return nil
}

func (d *DatabaseStore) GetThemeParkFromName(name string) (*ThemePark, error) {
	var themePark ThemePark
	err := d.conn.QueryRow(context.Background(),
		"select id, name, description, picture from themeparks where name = $1", name).
		Scan(&themePark.Id, &themePark.Name, &themePark.Description, &themePark.Picture)
	if err != nil {
		return nil, err
	}

	return &themePark, nil
}

func (d *DatabaseStore) AddThemePark(themePark ThemePark) error {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`insert into themeparks (name, location, description, picture) values ($1, POINT($2,$3), $4, $5)`,
		themePark.Name, themePark.Latitude, themePark.Longitude, themePark.Description, themePark.Picture)

	themeParkWithId, err := d.GetThemeParkFromName(themePark.Name)
	if err != nil {
		return err
	}

	// Insert new themeparks_categories
	for _, category := range themePark.Categories {
		_, err = tx.Exec(context.Background(),
			`insert into themeparks_categories (themepark_id, category_id) values ($1,$2)`, themeParkWithId.Id, category.Id)
	}

	tx.Commit(context.Background())
	return nil
}

func (d *DatabaseStore) DeleteThemePark(id int) error {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`DELETE FROM themeparks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())
	return nil
}

func (d *DatabaseStore) UpdateThemePark(themePark ThemePark) error {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`update themeparks set name = $1, location = POINT($2,$3), description = $4, picture = $5 where id = $6`,
		themePark.Name, themePark.Latitude, themePark.Longitude, themePark.Description, themePark.Picture, themePark.Id)
	if err != nil {
		return err
	}

	// Delete all categories for that themepark
	_, err = tx.Exec(context.Background(),
		`delete from themeparks_categories where themepark_id = $1`, themePark.Id)

	// Add new categories
	for _, category := range themePark.Categories {
		_, err = tx.Exec(context.Background(),
			`insert into themeparks_categories (themepark_id, category_id) values ($1,$2)`, themePark.Id, category.Id)
	}

	tx.Commit(context.Background())
	return nil
}
