package postgres

import (
	"coin/domain"
	"database/sql"
	"log"
	"time"
)

type preparedStatements struct {
	getUser       *sql.Stmt
	createUser    *sql.Stmt
	getOperations *sql.Stmt
}

type PostgresStore struct {
	DB    *sql.DB
	stmts *preparedStatements
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	InitCoinTables(db)
	stmts := newPreparedStatements(db)
	return &PostgresStore{
		DB:    db,
		stmts: stmts,
	}
}

func newPreparedStatements(db *sql.DB) *preparedStatements {
	const op = "database.postgres"
	stmts := &preparedStatements{}
	var err error
	stmts.getUser, err = db.Prepare(`
		SELECT u.id, u.username, u.balance, i.item_type, i.quantity 
		FROM users u
		LEFT JOIN inventory_items i on u.id=i.user_id
		WHERE u.username=$1;
	`)
	if err != nil {
		log.Println(op, err)
	}

	stmts.createUser, err = db.Prepare(`
		INSERT INTO users (username, balance)
		VALUES ($1, $2);
	`)
	if err != nil {
		log.Println(op, err)
	}

	stmts.getOperations, err = db.Prepare(`
		SELECT t.from_user, t.to_user, t.amount, t.created_at
		FROM transactions t
		JOIN users u1 ON t.from_user = u1.username
		JOIN users u2 ON t.to_user = u2.username
		WHERE t.from_user = $1 OR t.to_user = $1
		ORDER BY t.created_at DESC
	`)

	if err != nil {
		log.Println(op, err)
	}

	return stmts

}

func (r *PostgresStore) GetUser(username string) (domain.User, error) {
	var u domain.User
	var itemType sql.NullString
	var itemQuantity sql.NullInt64
	u.Inventory = make([]domain.InventoryItem, 0)
	rows, err := r.stmts.getUser.Query(username)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Username, &u.Balance, &itemType, &itemQuantity); err != nil {
			return domain.User{}, err
		}
		if itemType.Valid && itemQuantity.Valid {
			u.Inventory = append(u.Inventory, domain.InventoryItem{ItemType: itemType.String, Quantity: int(itemQuantity.Int64)})
		}
	}
	return u, nil
}

func (r *PostgresStore) CreateUser(username string, balance int) error {
	_, err := r.stmts.createUser.Exec(username, balance)
	return err
}

func (r *PostgresStore) PostBuyItem(userID uint, itemName string, price int) error {
	const op = "database.postgres.PostBuyItem"
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println(op, err)
		return err
	}

	_, err = tx.Exec(`
	UPDATE users SET balance = balance - $1 WHERE id = $2;`,
		price, userID)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Println(op, err2)
		}
		log.Println(op, err)
		return err
	}
	stmt := `INSERT INTO inventory_items (user_id, item_type, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, item_type)
		DO UPDATE SET quantity = inventory_items.quantity + EXCLUDED.quantity;`

	_, err = tx.Exec(stmt, userID, itemName)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Println(op, err2)
		}
		log.Println(op, err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Println(op, err)
		return err
	}
	return nil
}

func (r *PostgresStore) GetOperations(username string) ([]domain.Operations, error) {
	const op = "database.postgres.GetOperations"
	rows, err := r.stmts.getOperations.Query(username)
	operations := make([]domain.Operations, 0)
	if err != nil {
		return []domain.Operations{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var t domain.Operations
		err := rows.Scan(&t.FromUser, &t.ToUser, &t.Amount, &t.CreatedAt)
		if err != nil {
			log.Println(op, err)
			continue
		}
		operations = append(operations, t)
	}
	return operations, nil
}

func (r *PostgresStore) SendCoinTransaction(senderUsername, recipientUsername string, amount int) error {
	const op = "database.postgres.SendCoinTransaction"
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println(op, err)
		return err
	}

	_, err = tx.Exec(`
		UPDATE users SET balance = balance - $1 WHERE username = $2;`,
		amount, senderUsername)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Println(op, err2)
		}
		log.Println(op, err)
		return err
	}

	_, err = tx.Exec(`
		UPDATE users SET balance = balance + $1 WHERE username = $2;
	`, amount, recipientUsername)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Println(op, err2)
		}
		log.Println(op, err)
		return err
	}
	log.Printf("from_user = %s", senderUsername)
	_, err = tx.Exec(`INSERT INTO transactions (from_user, to_user, amount, created_at) VALUES ($1, $2, $3, $4)`,
		senderUsername, recipientUsername, amount, time.Now())
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Println(op, err2)
		}
		log.Println(op, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(op, err)
		return err
	}
	return nil
}
