package users

import (
	"context"
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	"github.com/monsterr00/gopher_mart/internal/helpers"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

type UserPostgresRepo struct {
	conn *sql.DB
}

func NewUserPostgresRepo(dbConn *sql.DB) *UserPostgresRepo {
	return &UserPostgresRepo{
		conn: dbConn,
	}
}

func (r *UserPostgresRepo) SaveUser(ctx context.Context, user entities.User) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCreateError)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users
		(ID, Login, Password)
		VALUES
		($1, $2, $3);
	    `, user.Id, user.Login, user.Password)
	if err != nil {
		return i.HandlePSQLError(err, i.MsgInsertError)
	}

	err = tx.Commit()
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCommitError)
	}

	return nil
}

func (r *UserPostgresRepo) GetUserByLogin(ctx context.Context, login string) (entities.User, error) {
	row := r.conn.QueryRowContext(ctx, `	
	SELECT 
		ID,
		Login,
		Password,
		CreatedAt,
		Balance,
		SpendedBonuses
	FROM users
	WHERE login = $1
    `,
		login,
	)

	var user entities.User
	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.CreatedAt, &user.Balance, &user.SpendedBonuses)
	if err != nil {
		return user, i.HandlePSQLError(err, i.MsgSelectError)
	}

	return user, nil
}

func (r *UserPostgresRepo) GetPassword(ctx context.Context, login string) (string, error) {
	row := r.conn.QueryRowContext(ctx, `	
	SELECT 		
		password
	FROM users
	WHERE login = $1
    `,
		login,
	)

	var password string
	err := row.Scan(&password)
	if err != nil {
		return password, i.HandlePSQLError(err, i.MsgSelectError)
	}

	return password, nil
}

func (r *UserPostgresRepo) CheckDuplicateLogin(ctx context.Context, login string) error {
	row := r.conn.QueryRowContext(ctx, `	
	SELECT 		
		id
	FROM users
	WHERE login = $1
    `,
		login,
	)

	var id string
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	}

	return i.ErrAlreadyExists
}
