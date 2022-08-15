package usecase

import (
	"github.com/mihailov-vf/codebank/domain"
	"github.com/mihailov-vf/codebank/dto"
)

type UsecaseTransaction struct {
	transactionRepository domain.TransactionRepository
}

func NewUsecaseTransaction(transactionRepository domain.TransactionRepository) *UsecaseTransaction {
	return &UsecaseTransaction{transactionRepository}
}

func (u UsecaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	creditCard := u.hydrateCreditCard(transactionDto)
	creditCard, err := u.transactionRepository.GetCreditCard(creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}
	t := u.NewTransaction(transactionDto, creditCard)
	t.ProcessAndValidate(&creditCard)
	err = u.transactionRepository.SaveTransaction(*t, creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}
	return *t, nil
}

func (u UsecaseTransaction) hydrateCreditCard(transactionDto dto.Transaction) domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDto.Name
	creditCard.Number = transactionDto.Number
	creditCard.ExpirationMonth = transactionDto.ExpirationMonth
	creditCard.ExpirationYear = transactionDto.ExpirationYear
	creditCard.CVV = transactionDto.CVV
	return *creditCard
}

func (u UsecaseTransaction) NewTransaction(transactiondto dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()
	t.CreditCardId = cc.ID
	t.Amount = transactiondto.Amount
	t.Description = transactiondto.Description
	t.Store = transactiondto.Store
	return t
}
