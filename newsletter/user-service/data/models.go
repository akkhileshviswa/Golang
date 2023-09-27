package data

import (
	"context"
	"database/sql"
	"time"
)

const DbTimeout = time.Second * 3

var Db *sql.DB

func New(DbPool *sql.DB) Models {
	Db = DbPool

	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Is_sent   bool      `json:"is_sent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessagePayload struct {
	ID        int
	FirstName string `bson:"first_name" json:"first_name"`
	Name      string `bson:"name" json:"name"`
	To        string `bson:"to" json:"to"`
	Subject   string `bson:"subject" json:"subject"`
	Message   string `bson:"message" json:"message"`
}

// This function is used to update is_mail sent field if mail has been sent successfully
func UpdateIfMailSent(payload MessagePayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	stmt := `update users set
		is_sent = $1,
		updated_at = $2
		where id = $3
	`

	_, err := Db.ExecContext(ctx, stmt,
		true,
		time.Now(),
		payload.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
