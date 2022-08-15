package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mihailov-vf/codebank/domain"
	"github.com/mihailov-vf/codebank/infrastructure/repository"
	"github.com/mihailov-vf/codebank/usecase"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Wesley"
	cc.ExpirationYear = 2021
	cc.ExpirationMonth = 7
	cc.CVV = 123
	cc.Limit = 1000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryPostgres(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}
}

func setupTransactionUseCase(db *sql.DB) usecase.UsecaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryPostgres(db)
	return *usecase.NewUsecaseTransaction(transactionRepository)
}

func setupDb() *sql.DB {
	pgsqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"db",
		5432,
		"root",
		"root",
		"codebank")
	db, err := sql.Open("postgres", pgsqlInfo)
	if err != nil {
		log.Fatal("Database connection error")
	}
	return db
}
