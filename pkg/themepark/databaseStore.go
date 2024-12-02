package themepark

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (d *DatabaseStore) IsLoggedIn(token string) bool {
	rows, _ := d.conn.Query(context.Background(),
		"select * from tokens where token = $1", token)

	for rows.Next() { // This means a register exist, meaning the user is logged in
		rows.Close()
		return true
	}

	return false
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
