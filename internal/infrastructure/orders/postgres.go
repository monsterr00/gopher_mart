package orders

import (
	"context"
	"database/sql"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	"github.com/monsterr00/gopher_mart/internal/helpers"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

type OrderPostgresRepo struct {
	conn *sql.DB
}

func NewOrderPostgresRepo(dbConn *sql.DB) *OrderPostgresRepo {
	return &OrderPostgresRepo{
		conn: dbConn,
	}
}

func (r *OrderPostgresRepo) GetOrdersLogin(ctx context.Context, login string) ([]entities.Order, error) {
	rows, err := r.conn.QueryContext(ctx, `
	SELECT		
		OrderNum,		
		CreatedAt,
		AddedBonuses,
		Status
	FROM orders
	WHERE login = $1
	ORDER BY CreatedAt ASC
    `,
		login,
	)

	if err != nil {
		return nil, helpers.WrapError(err, i.MsgBuilderSelect)
	}

	defer rows.Close()

	var orders []entities.Order

	for rows.Next() {
		var order entities.Order

		err := rows.Scan(&order.OrderNum, &order.CreatedAt, &order.AddedBonuses, &order.Status)
		if err != nil {
			return nil, i.HandlePSQLError(err, i.MsgSelectError)
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, helpers.WrapError(err, i.MsgIterationError)
	}

	if len(orders) == 0 {
		return nil, i.ErrNoData
	}

	return orders, nil
}

func (r *OrderPostgresRepo) SaveOrder(ctx context.Context, order entities.Order) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCreateError)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders
		(ID, OrderNum, Login)
		VALUES
		($1, $2, $3);
	    `, order.Id, order.OrderNum, order.Login)
	if err != nil {
		return i.HandlePSQLError(err, i.MsgInsertError)
	}

	err = tx.Commit()
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCommitError)
	}

	return nil
}

func (r *OrderPostgresRepo) CheckOrderDuplicate(ctx context.Context, order entities.Order) error {
	row := r.conn.QueryRowContext(ctx, `	
	SELECT 
		ordernum		
	FROM orders
	WHERE login = $1 and
	      ordernum = $2
    `,
		order.Login, order.OrderNum,
	)

	var ordernum string
	err := row.Scan(&ordernum)
	if err != nil && err != sql.ErrNoRows {
		return i.HandlePSQLError(err, i.MsgSelectError)
	}

	if ordernum != "" {
		return i.ErrDuplicateUserOrder
	}

	row = r.conn.QueryRowContext(ctx, `	
	SELECT 
		ID		
	FROM orders
	WHERE ordernum = $1
    `,
		order.OrderNum,
	)

	err = row.Scan(&ordernum)
	if err != nil && err != sql.ErrNoRows {
		return i.HandlePSQLError(err, i.MsgSelectError)
	}

	if ordernum != "" {
		return i.ErrDuplicateOrder
	}

	return nil
}

func (r *OrderPostgresRepo) UpdateStatus(ctx context.Context, order entities.RegistredOrder) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCreateError)
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE orders
	SET Status = $1
	WHERE OrderNum = $2;		
    `, order.Status, order.OrderNum)
	if err != nil {
		return i.HandlePSQLError(err, i.MsgUpdateError)
	}

	err = tx.Commit()
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCommitError)
	}

	return nil
}

func (r *OrderPostgresRepo) UpdateBonuses(ctx context.Context, order entities.RegistredOrder) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCreateError)
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE orders
	SET AddedBonuses = $1
	WHERE OrderNum = $2;		
    `, order.AddedBonuses, order.OrderNum)
	if err != nil {
		return i.HandlePSQLError(err, i.MsgUpdateError)
	}

	err = tx.Commit()
	if err != nil {
		return helpers.WrapError(err, i.MsgTXCommitError)
	}

	return nil
}
