package user

import (
	"log"
	"shift-be/src/lib"

	sq "github.com/Masterminds/squirrel"
)

func SeedUsers(psql sq.StatementBuilderType) {
	password, _ := lib.HashPassword("password")

	fakeUsers := []struct {
		Fullname string
		Username string
		Email    string
		Password string
	}{
		{"Daffa Abdullah", "daffa.abdllh", "daffa@example.com", password},
	}

	query := psql.Insert("users").Columns("fullname", "username", "email", "password")
	for _, u := range fakeUsers {
		query = query.Values(u.Fullname, u.Username, u.Email, u.Password)
	}

	_, err := query.Exec()
	if err != nil {
		log.Printf("Failed to seed users: %v", err)
		return
	}

	log.Println("Users seeded successfully!")
}
