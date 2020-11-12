package commands

import (
	"fmt"

	"github.com/dogmatiq/example/messages"
)

// OpenAccountForNewCustomer is a command requesting that a new bank account be
// opened for a new customer.
type OpenAccountForNewCustomer struct {
	CustomerID   string
	CustomerName string
	AccountID    string
	AccountName  string
}

// OpenAccount is a command requesting that a new bank account be opened for an
// existing customer.
type OpenAccount struct {
	CustomerID  string
	AccountID   string
	AccountName string
}

// CreditAccount is a command that requests a bank account be credited.
type CreditAccount struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
}

// DebitAccount is a command that requests a bank account be debited.
type DebitAccount struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
	ScheduledDate   string
}

// MessageDescription returns a human-readable description of the message.
func (m OpenAccountForNewCustomer) MessageDescription() string {
	return fmt.Sprintf(
		"customer %s %s is opening their first account %s %s",
		m.CustomerID,
		m.CustomerName,
		m.AccountID,
		m.AccountName,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m OpenAccount) MessageDescription() string {
	return fmt.Sprintf(
		"opening account %s %s for customer %s",
		m.AccountID,
		m.AccountName,
		m.CustomerID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m CreditAccount) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: crediting %s to account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DebitAccount) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: debiting %s from account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}
