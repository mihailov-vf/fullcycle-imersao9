package repository

import (
	"database/sql"
	"errors"

	"github.com/mihailov-vf/codebank/domain"
)

type TransactionRepositoryPostgres struct {
	db *sql.DB
}

func NewTransactionRepositoryPostgres(db *sql.DB) *TransactionRepositoryPostgres {
	return &TransactionRepositoryPostgres{db: db}
}

func (t TransactionRepositoryPostgres) SaveTransaction(transaction domain.Transaction, creditcard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`INSERT INTO 
	transactions(id, credit_card_id, amount, status, description, store, created_at)
	VALUES($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		transaction.ID,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}
	if transaction.Status != "rejected" {
		err = t.updateBalance(creditcard)
		if err != nil {
			return err
		}
	}
	return stmt.Close()
}

func (t TransactionRepositoryPostgres) updateBalance(creditcard domain.CreditCard) error {
	_, err := t.db.Exec("UPDATE credit_cards set balance = $1 WHERE id = $2", creditcard.Balance, creditcard.ID)
	return err
}

func (t TransactionRepositoryPostgres) CreateCreditCard(creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`INSERT INTO credit_cards(id, name, number, expiration_month, expiration_year, cvv, balance, balance_limit)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (t TransactionRepositoryPostgres) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare("select id, balance, balance_limit from credit_cards where number=$1")
	if err != nil {
		return c, err
	}
	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}
