package service

import (
	"math"

	"github.com/inflop/splitty/internal/domain/model"
)

// ExpenseService obsługuje operacje na wydatkach i rozliczeniach
type ExpenseService struct{}

// NewExpenseService tworzy nową instancję usługi wydatków
func NewExpenseService() *ExpenseService {
	return &ExpenseService{}
}

// RoundToTwo zaokrągla liczbę do dwóch miejsc po przecinku
func (s *ExpenseService) RoundToTwo(num float64) float64 {
	return math.Round(num*100) / 100
}

// CalculateSummary oblicza podsumowanie wydarzenia
func (s *ExpenseService) CalculateSummary(event *model.Event) *model.Summary {
	// Przetwarzanie wydatków
	totalAmount := 0.0

	type processedExpense struct {
		expense              model.Expense
		processedTotalAmount float64
		processedPayments    float64
	}

	processedExpenses := make([]processedExpense, len(event.Expenses))

	for i, exp := range event.Expenses {
		// Suma wszystkich płatności
		paymentsSum := 0.0
		for _, payment := range exp.Payments {
			paymentsSum = s.RoundToTwo(paymentsSum + payment.Amount)
		}

		// Jeśli podano totalAmount, używamy go, w przeciwnym razie suma płatności
		totalExp := exp.TotalAmount
		if totalExp == 0 {
			totalExp = paymentsSum
		}

		processedExpenses[i] = processedExpense{
			expense:              exp,
			processedTotalAmount: totalExp,
			processedPayments:    paymentsSum,
		}

		totalAmount = s.RoundToTwo(totalAmount + totalExp)
	}

	// Średnia na osobę
	perPersonAmount := 0.0
	if len(event.Participants) > 0 {
		perPersonAmount = s.RoundToTwo(totalAmount / float64(len(event.Participants)))
	}

	// Ile każdy zapłacił
	paidByPerson := make([]model.ParticipantBalance, len(event.Participants))

	for i, person := range event.Participants {
		paidAmount := 0.0

		// Obliczanie ile osoba zapłaciła
		for _, exp := range event.Expenses {
			for _, payment := range exp.Payments {
				if payment.ParticipantID == person.ID {
					paidAmount = s.RoundToTwo(paidAmount + payment.Amount)
				}
			}
		}

		// Obliczanie ile osoba powinna zapłacić
		shouldPay := 0.0
		for _, exp := range processedExpenses {
			isShared := false
			for _, id := range exp.expense.SharedWith {
				if id == person.ID {
					isShared = true
					break
				}
			}

			if isShared {
				sharedCount := len(exp.expense.SharedWith)
				if sharedCount == 0 {
					sharedCount = 1
				}
				perPersonInExpense := s.RoundToTwo(exp.processedTotalAmount / float64(sharedCount))
				shouldPay = s.RoundToTwo(shouldPay + perPersonInExpense)
			}
		}

		// Bilans
		balance := s.RoundToTwo(paidAmount - shouldPay)

		paidByPerson[i] = model.ParticipantBalance{
			ID:        person.ID,
			Name:      person.Name,
			Paid:      paidAmount,
			ShouldPay: shouldPay,
			Balance:   balance,
		}
	}

	// Obliczanie rozliczeń
	settlements := s.CalculateSettlements(paidByPerson)

	return &model.Summary{
		TotalAmount:     totalAmount,
		PerPersonAmount: perPersonAmount,
		PaidByPerson:    paidByPerson,
		Settlements:     settlements,
	}
}

// CalculateSettlements oblicza rozliczenia między uczestnikami
func (s *ExpenseService) CalculateSettlements(balances []model.ParticipantBalance) []model.Settlement {
	// Kopiowanie do struktur roboczych
	type workBalance struct {
		balance          model.ParticipantBalance
		remainingBalance float64
	}

	// Identyfikacja dłużników (balans ujemny)
	var debtorsWork []workBalance
	for _, b := range balances {
		if b.Balance < -0.02 {
			debtorsWork = append(debtorsWork, workBalance{
				balance:          b,
				remainingBalance: b.Balance,
			})
		}
	}

	// Identyfikacja wierzycieli (balans dodatni)
	var creditorsWork []workBalance
	for _, b := range balances {
		if b.Balance > 0.02 {
			creditorsWork = append(creditorsWork, workBalance{
				balance:          b,
				remainingBalance: b.Balance,
			})
		}
	}

	// Sortowanie dłużników (rosnąco po bilansie)
	for i := 0; i < len(debtorsWork)-1; i++ {
		for j := i + 1; j < len(debtorsWork); j++ {
			if debtorsWork[i].remainingBalance > debtorsWork[j].remainingBalance {
				debtorsWork[i], debtorsWork[j] = debtorsWork[j], debtorsWork[i]
			}
		}
	}

	// Sortowanie wierzycieli (malejąco po bilansie)
	for i := 0; i < len(creditorsWork)-1; i++ {
		for j := i + 1; j < len(creditorsWork); j++ {
			if creditorsWork[i].remainingBalance < creditorsWork[j].remainingBalance {
				creditorsWork[i], creditorsWork[j] = creditorsWork[j], creditorsWork[i]
			}
		}
	}

	var settlements []model.Settlement
	debtIndex := 0
	creditIndex := 0

	// Algorytm rozliczeń
	for debtIndex < len(debtorsWork) && creditIndex < len(creditorsWork) {
		debtor := &debtorsWork[debtIndex]
		creditor := &creditorsWork[creditIndex]

		amount := s.RoundToTwo(math.Min(math.Abs(debtor.remainingBalance), creditor.remainingBalance))

		if amount > 0.02 {
			settlements = append(settlements, model.Settlement{
				From:     debtor.balance.ID,
				FromName: debtor.balance.Name,
				To:       creditor.balance.ID,
				ToName:   creditor.balance.Name,
				Amount:   amount,
			})
		}

		debtor.remainingBalance = s.RoundToTwo(debtor.remainingBalance + amount)
		creditor.remainingBalance = s.RoundToTwo(creditor.remainingBalance - amount)

		if math.Abs(debtor.remainingBalance) < 0.02 {
			debtIndex++
		}
		if creditor.remainingBalance < 0.02 {
			creditIndex++
		}
	}

	return settlements
}
