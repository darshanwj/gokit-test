package model

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Id    uint        `json:"id"`
	Name  string      `json:"name"`
	Phone null.String `json:"phone"`
}

type userRepository struct {
	db *sql.DB
}

type UserRepository interface {
	FindOneById(id uint) (User, error)
	FindMany() []User
	Insert(name, phone string) int64
}

func NewUserRespository(db *sql.DB) UserRepository {
	return userRepository{db: db}
}

func (r userRepository) FindOneById(id uint) (User, error) {
	var user User
	err := sq.
		Select("*").
		From("user").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).Scan(&user.Id, &user.Name, &user.Phone)

	if err != nil {
		log.Print(err)
		return User{}, err
	}

	return user, nil
}

func (r userRepository) FindMany() []User {
	res, err := sq.
		Select("*").
		From("user").
		RunWith(r.db).Query()

	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	var users []User
	for res.Next() {
		var user User
		err = res.Scan(&user.Id, &user.Name, &user.Phone)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}

func (r userRepository) Insert(name, phone string) int64 {
	res, err := sq.
		Insert("user").
		Columns("name", "phone").
		Values(name, phone).
		RunWith(r.db).Exec()

	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return id
}
