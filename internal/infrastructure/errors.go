package infrastructure

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/monsterr00/gopher_mart/internal/helpers"
)

var (
	ErrAlreadyExists      = errors.New("record already exists")
	ErrDoesNotExist       = errors.New("record does not exist")
	ErrDuplicateUserOrder = errors.New("duplicate order for current user")
	ErrDuplicateOrder     = errors.New("duplicate order")
	ErrNoData             = errors.New("no data")
	ErrValidate           = errors.New("validation error")
	ErrInsufFunds         = errors.New("insufficient funds")
)

// error messages
const (
	MsgSelectError    = "select error"
	MsgInsertError    = "insert error"
	MsgUpdateError    = "update error"
	MsgBuilderSelect  = "can't create select query"
	MsgIterationError = "select result iteration error"
	MsgTXCreateError  = "can't create transaction"
	MsgTXCommitError  = "can't commit transaction"
)

func HandlePSQLError(err error, description string) error {
	if err == sql.ErrNoRows {
		return ErrDoesNotExist
	}

	pqErr, ok := err.(*pq.Error)
	if ok {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return ErrAlreadyExists
		}
	}

	return helpers.WrapError(err, description)
}
