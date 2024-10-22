package db

import (
	"log"
	"video-conference/pkg/models"
)

func (repo *RepositoryImpl) SaveUser(user *models.User) error {
	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err := repo.DB.Exec(query, user.ID, user.Username, user.Email, user.Password)
	return err
}

func (repo *RepositoryImpl) FindUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	row := repo.DB.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *RepositoryImpl) FindUsersByUsername(name string) ([]models.User, error) {
	query := `SELECT id, username, email FROM users WHERE username ILIKE '%' || $1 || '%' LIMIT 5`
	rows, err := repo.DB.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *RepositoryImpl) FindUserById(id string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE id = $1`

	row := repo.DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}

	return &user, nil
}
