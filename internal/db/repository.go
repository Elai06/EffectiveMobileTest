package db

import (
	"database/sql"
	"effectiveMobile/env"
	"effectiveMobile/internal"
	"effectiveMobile/internal/utils"
	"fmt"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg env.Config) (*Repository, error) {
	dbUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSlMode,
	)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	utils.DebugLog("connected to database", dbUrl)

	return &Repository{db}, nil
}

func (r *Repository) CreatePerson(person internal.Person) error {
	query := "INSERT INTO people (name, surname, age, gender, nationality) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.Exec(query, person.Name, person.Surname, person.Age, person.Gender, person.Nationality)
	if err != nil {
		return fmt.Errorf("error creating person: %w", err)
	}

	utils.DebugLog("person created", person)
	return nil
}

func (r *Repository) GetPeople(personFilter internal.PersonFilter) ([]internal.Person, error) {
	query :=
		`SELECT name, surname, age, gender, nationality, id FROM people
        WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
        AND ($2 = '' OR surname ILIKE '%' || $2 || '%')
        AND ($3 = '' OR gender ILIKE '%' || $3 || '%')
        AND ($5 = '' OR nationality ILIKE '%' || $5 || '%')
        AND (CASE WHEN $4 > 0 THEN age = $4 ELSE TRUE END)
        AND (CASE WHEN $6 > 0 THEN age >= $6 ELSE TRUE END)
        AND (CASE WHEN $7 > 0 THEN age <= $7 ELSE TRUE END)
        ORDER BY id DESC
        LIMIT $8 OFFSET $9`

	rows, err := r.db.Query(query,
		personFilter.Person.Name,
		personFilter.Person.Surname,
		personFilter.Person.Gender,
		personFilter.Person.Age,
		personFilter.Person.Nationality,
		personFilter.AgeMin,
		personFilter.AgeMax,
		personFilter.Limit,
		personFilter.Offset)
	if err != nil {
		return nil, fmt.Errorf("error creating person: %w", err)
	}

	var people []internal.Person
	for rows.Next() {
		person := internal.Person{}
		err = rows.Scan(&person.Name, &person.Surname, &person.Age, &person.Gender, &person.Nationality, &person.ID)
		if err != nil {
			return nil, fmt.Errorf("error creating person: %w", err)
		}
		people = append(people, person)
	}

	utils.DebugLog("get people", people)

	return people, nil
}

func (r *Repository) UpdatePerson(person internal.Person) error {
	query := `UPDATE people SET 
                  name = $2,
                  surname = $3,
                  age = $4, 
                  gender = $5,
                  nationality = $6 
               WHERE id = $1`

	res, err := r.db.Exec(query, person.ID, person.Name, person.Surname, person.Age, person.Gender, person.Nationality)
	if err != nil {
		return fmt.Errorf("error creating person: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no person updated")

	}

	utils.DebugLog("person updated", person)
	return nil
}

func (r *Repository) DeletePerson(id int) error {
	query := "DELETE FROM people WHERE id = $1"

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}

	utils.DebugLog("person deleted", id)

	return nil
}
