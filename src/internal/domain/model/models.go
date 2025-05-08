package model

// Participant reprezentuje uczestnika wydarzenia
type Participant struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

// Payment reprezentuje pojedynczą płatność w ramach wydatku
type Payment struct {
	ParticipantID int     `json:"participantId"`
	Amount        float64 `json:"amount"`
}

// Expense reprezentuje wydatek grupowy
type Expense struct {
	ID          int       `json:"id"`
	Category    string    `json:"category"`
	TotalAmount float64   `json:"totalAmount"`
	Payments    []Payment `json:"payments"`
	SharedWith  []int     `json:"sharedWith"`
}

// Event reprezentuje całe wydarzenie z uczestnikami i wydatkami
type Event struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Participants []Participant `json:"participants"`
	Expenses     []Expense     `json:"expenses"`
}

// ParticipantBalance zawiera informacje o bilansie uczestnika
type ParticipantBalance struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Paid      float64 `json:"paid"`
	ShouldPay float64 `json:"shouldPay"`
	Balance   float64 `json:"balance"`
}

// Settlement reprezentuje pojedyncze rozliczenie między uczestnikami
type Settlement struct {
	From     int     `json:"from"`
	FromName string  `json:"fromName"`
	To       int     `json:"to"`
	ToName   string  `json:"toName"`
	Amount   float64 `json:"amount"`
}

// Summary reprezentuje podsumowanie wydarzenia
type Summary struct {
	TotalAmount     float64              `json:"totalAmount"`
	PerPersonAmount float64              `json:"perPersonAmount"`
	PaidByPerson    []ParticipantBalance `json:"paidByPerson"`
	Settlements     []Settlement         `json:"settlements"`
}
