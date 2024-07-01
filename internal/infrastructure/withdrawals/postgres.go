package withdrawals

import (
	"context"
	"database/sql"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	"github.com/monsterr00/gopher_mart/internal/helpers"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

type WithdrawalPostgresRepo struct {
	conn *sql.DB
}

func NewWithdrawalPostgresRepo(dbConn *sql.DB) *WithdrawalPostgresRepo {
	return &WithdrawalPostgresRepo{
		conn: dbConn,
	}
}

func (w WithdrawalPostgresRepo) GetBalanceAdded(ctx context.Context, login string) (float64, error) {
	row := w.conn.QueryRowContext(ctx, `	
	SELECT 								
		SUM (AddedBonuses)		
	FROM orders
	WHERE login = $1 
	AND AddedBonuses is not null
	GROUP BY login
    `,
		login,
	)

	var bonuses float64
	err := row.Scan(&bonuses)
	if err != nil && err != sql.ErrNoRows {
		return bonuses, i.HandlePSQLError(err, i.MsgSelectError)
	}

	return bonuses, nil
}

func (w WithdrawalPostgresRepo) GetBalanceSpended(ctx context.Context, login string) (float64, error) {
	row := w.conn.QueryRowContext(ctx, `	
	SELECT 								
		SUM (SpendedBonuses)		
	FROM orders
	WHERE login = $1 
	AND SpendedBonuses is not null
	GROUP BY login
    `,
		login,
	)

	var bonuses float64
	err := row.Scan(&bonuses)
	if err != nil && err != sql.ErrNoRows {
		return bonuses, i.HandlePSQLError(err, i.MsgSelectError)
	}

	return bonuses, nil
}

func (w WithdrawalPostgresRepo) GetWithdrawals(ctx context.Context, login string) ([]entities.Withdrawal, error) {
	rows, err := w.conn.QueryContext(ctx, `
	SELECT		
		OrderNum,				
		SpendedBonuses,
		SpendDate
	FROM orders
	WHERE login = $1
	and SpendedBonuses > 0
	ORDER BY SpendDate ASC
    `,
		login,
	)

	if err != nil {
		return nil, helpers.WrapError(err, i.MsgBuilderSelect)
	}

	defer rows.Close()

	var withdrawals []entities.Withdrawal

	for rows.Next() {
		var withdrawal entities.Withdrawal

		err := rows.Scan(&withdrawal.OrderNum, &withdrawal.Bonuses, &withdrawal.SpendDate)
		if err != nil {
			return nil, i.HandlePSQLError(err, i.MsgSelectError)
		}

		withdrawals = append(withdrawals, withdrawal)
	}
	if err := rows.Err(); err != nil {
		return nil, helpers.WrapError(err, i.MsgIterationError)
	}

	if len(withdrawals) == 0 {
		return nil, i.ErrNoData
	}

	return withdrawals, nil
}

func (w WithdrawalPostgresRepo) SaveWithdrawal(ctx context.Context, order entities.Order) error {
	tx, err := w.conn.BeginTx(ctx, nil)
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCreateError)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE orders
		SET SpendedBonuses = $1, 
		SpendDate = now()
		WHERE orderNUm = $2;
	    `, order.SpendedBonuses, order.OrderNum)
	if err != nil {
		return i.HandlePSQLError(err, i.MsgInsertError)
	}

	err = tx.Commit()
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCommitError)
	}

	return nil
}
