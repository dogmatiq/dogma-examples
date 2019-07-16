package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// transfer is the process root for a funds transfer.
type transferProcess struct {
	FromAccountID string
	ToAccountID   string
}

// TransferProcessHandler manages the process of transferring funds between
// accounts.
type TransferProcessHandler struct {
	dogma.NoTimeoutBehavior
}

// New returns a new transfer instance.
func (TransferProcessHandler) New() dogma.ProcessRoot {
	return &transferProcess{}
}

// Configure configures the behavior of the engine as it relates to this handler.
func (TransferProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("transfer")

	c.ConsumesEventType(events.TransferStarted{})
	c.ConsumesEventType(events.AccountDebited{})
	c.ConsumesEventType(events.AccountDebitDeclined{})
	c.ConsumesEventType(events.DailyDebitLimitConsumed{})
	c.ConsumesEventType(events.DailyDebitLimitExceeded{})
	c.ConsumesEventType(events.AccountCredited{})
	c.ConsumesEventType(events.TransferApproved{})
	c.ConsumesEventType(events.TransferDeclined{})

	c.ProducesCommandType(commands.DebitAccount{})
	c.ProducesCommandType(commands.ConsumeDailyDebitLimit{})
	c.ProducesCommandType(commands.CreditAccount{})
	c.ProducesCommandType(commands.ApproveTransfer{})
	c.ProducesCommandType(commands.DeclineTransfer{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (TransferProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.TransferStarted:
		return x.TransactionID, true, nil
	case events.AccountDebited:
		return x.TransactionID, x.TransactionType == messages.Transfer, nil
	case events.AccountDebitDeclined:
		return x.TransactionID, x.TransactionType == messages.Transfer, nil
	case events.DailyDebitLimitConsumed:
		return x.TransactionID, x.DebitType == messages.Transfer, nil
	case events.DailyDebitLimitExceeded:
		return x.TransactionID, x.DebitType == messages.Transfer, nil
	case events.AccountCredited:
		return x.TransactionID, x.TransactionType == messages.Transfer, nil
	case events.TransferApproved:
		return x.TransactionID, true, nil
	case events.TransferDeclined:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (TransferProcessHandler) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.TransferStarted:
		s.Begin()

		r := s.Root().(*transferProcess)
		r.FromAccountID = x.FromAccountID
		r.ToAccountID = x.ToAccountID

		s.ExecuteCommand(commands.DebitAccount{
			TransactionID:   x.TransactionID,
			AccountID:       x.FromAccountID,
			TransactionType: messages.Transfer,
			Amount:          x.Amount,
			ScheduledDate:   x.ScheduledDate,
		})

	case events.AccountDebited:
		s.ExecuteCommand(commands.ConsumeDailyDebitLimit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			DebitType:     messages.Transfer,
			Amount:        x.Amount,
			ScheduledDate: x.ScheduledDate,
		})

	case events.AccountDebitDeclined:
		r := s.Root().(*transferProcess)

		s.ExecuteCommand(commands.DeclineTransfer{
			TransactionID: x.TransactionID,
			FromAccountID: r.FromAccountID,
			ToAccountID:   r.ToAccountID,
			Amount:        x.Amount,
			Reason:        x.Reason,
		})

	case events.DailyDebitLimitConsumed:
		r := s.Root().(*transferProcess)

		// continue transfer
		s.ExecuteCommand(commands.CreditAccount{
			TransactionID:   x.TransactionID,
			AccountID:       r.ToAccountID,
			TransactionType: messages.Transfer,
			Amount:          x.Amount,
		})

	case events.DailyDebitLimitExceeded:
		r := s.Root().(*transferProcess)

		// compensate the initial debit
		s.ExecuteCommand(commands.CreditAccount{
			TransactionID:   x.TransactionID,
			AccountID:       r.FromAccountID,
			TransactionType: messages.Transfer,
			Amount:          x.Amount,
		})

		s.ExecuteCommand(commands.DeclineTransfer{
			TransactionID: x.TransactionID,
			FromAccountID: r.FromAccountID,
			ToAccountID:   r.ToAccountID,
			Amount:        x.Amount,
			Reason:        messages.DailyDebitLimitExceeded,
		})

	case events.AccountCredited:
		r := s.Root().(*transferProcess)

		// check if it was a credit to complete the transfer (success) and not
		// a compensation credit (failure)
		if r.ToAccountID == x.AccountID {
			s.ExecuteCommand(commands.ApproveTransfer{
				TransactionID: x.TransactionID,
				FromAccountID: r.FromAccountID,
				ToAccountID:   r.ToAccountID,
				Amount:        x.Amount,
			})
		}

	case events.TransferApproved,
		events.TransferDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
