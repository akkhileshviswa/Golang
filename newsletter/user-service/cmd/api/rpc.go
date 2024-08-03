package main

import (
	"context"
	"errors"
	"log"
	"time"

	"usermanagement/data"
)

type SubscribePayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Response struct {
	Result []data.User
}

// Insert Record writes our payload to postgres.
func (r *RPCServer) InsertRecord(s SubscribePayload, response *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), data.DbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into users (first_name, last_name, email, is_sent, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	if !ValidateName(s.FirstName) || !ValidateName(s.LastName) || !ValidateEmail(s.Email) {
		return errors.New("enter valid details")
	}

	err := data.Db.QueryRowContext(ctx, stmt,
		s.FirstName,
		s.LastName,
		s.Email,
		false,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	*response = newID
	if err != nil {
		return err
	}

	return nil
}

// Will get the user records based on is_sent status.
func (r *RPCServer) GetMailUserRecords(_ struct{}, resp *Response) error {
	log.Println("Inside user")
	ctx, cancel := context.WithTimeout(context.Background(), data.DbTimeout)
	defer cancel()

	query := `select id, first_name, last_name, email, is_sent, created_at, updated_at
	from users where is_sent = $1`

	rows, err := data.Db.QueryContext(ctx, query, false)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	var users []data.User

	for rows.Next() {
		var user data.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Is_sent,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	resp.Result = users

	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	return nil
}
