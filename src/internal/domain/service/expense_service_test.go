package service_test

import (
	"testing"

	"github.com/inflop/splitty.api/internal/domain/model"
	"github.com/inflop/splitty.api/internal/domain/service"
)

func TestCalculateSummary(t *testing.T) {
	// Przygotowanie danych testowych
	event := &model.Event{
		ID:   1,
		Name: "Test Event",
		Participants: []model.Participant{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Charlie"},
		},
		Expenses: []model.Expense{
			{
				ID:          1,
				Category:    "Accommodation",
				TotalAmount: 300,
				Payments: []model.Payment{
					{ParticipantID: 1, Amount: 300},
				},
				SharedWith: []int{1, 2, 3},
			},
			{
				ID:          2,
				Category:    "Food",
				TotalAmount: 150,
				Payments: []model.Payment{
					{ParticipantID: 2, Amount: 150},
				},
				SharedWith: []int{1, 2, 3},
			},
			{
				ID:          3,
				Category:    "Transport",
				TotalAmount: 90,
				Payments: []model.Payment{
					{ParticipantID: 3, Amount: 90},
				},
				SharedWith: []int{1, 2, 3},
			},
		},
	}

	// Utworzenie usługi
	expenseService := service.NewExpenseService()

	// Wywołanie metody
	summary := expenseService.CalculateSummary(event)

	// Sprawdzenie wyników
	if summary.TotalAmount != 540 {
		t.Errorf("Expected total amount to be 540, got %v", summary.TotalAmount)
	}

	if summary.PerPersonAmount != 180 {
		t.Errorf("Expected per person amount to be 180, got %v", summary.PerPersonAmount)
	}

	// Sprawdzenie bilansu każdego uczestnika
	for _, balance := range summary.PaidByPerson {
		switch balance.ID {
		case 1: // Alice zapłaciła 300, powinna zapłacić 180, bilans +120
			if balance.Balance != 120 {
				t.Errorf("Expected Alice's balance to be 120, got %v", balance.Balance)
			}
		case 2: // Bob zapłacił 150, powinien zapłacić 180, bilans -30
			if balance.Balance != -30 {
				t.Errorf("Expected Bob's balance to be -30, got %v", balance.Balance)
			}
		case 3: // Charlie zapłacił 90, powinien zapłacić 180, bilans -90
			if balance.Balance != -90 {
				t.Errorf("Expected Charlie's balance to be -90, got %v", balance.Balance)
			}
		}
	}

	// Sprawdzenie rozliczeń
	if len(summary.Settlements) != 2 {
		t.Errorf("Expected 2 settlements, got %d", len(summary.Settlements))
	}

	// Sprawdzenie szczegółów rozliczeń
	for _, settlement := range summary.Settlements {
		if settlement.From == 2 && settlement.To == 1 {
			if settlement.Amount != 30 {
				t.Errorf("Expected Bob to pay Alice 30, got %v", settlement.Amount)
			}
		} else if settlement.From == 3 && settlement.To == 1 {
			if settlement.Amount != 90 {
				t.Errorf("Expected Charlie to pay Alice 90, got %v", settlement.Amount)
			}
		} else {
			t.Errorf("Unexpected settlement from %d to %d", settlement.From, settlement.To)
		}
	}
}
