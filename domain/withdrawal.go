package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// WithdrawalProcess manages the process of withdrawing funds from an account.
type WithdrawalProcess struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (WithdrawalProcess) Configure(c dogma.ProcessConfigurer) {
	c.Name("withdrawal")

	c.ConsumesEventType(events.WithdrawalStarted{})
	c.ConsumesEventType(events.FundsHeldForWithdrawal{})
	c.ConsumesEventType(events.WithdrawalDeclined{})
	c.ConsumesEventType(events.DailyDebitLimitConsumed{})
	c.ConsumesEventType(events.DailyDebitLimitExceeded{})
	c.ConsumesEventType(events.AccountDebitedForWithdrawal{})

	c.ProducesCommandType(commands.ConsumeDailyDebitLimit{})
	c.ProducesCommandType(commands.DeclineWithdrawal{})
	c.ProducesCommandType(commands.HoldFundsForWithdrawal{})
	c.ProducesCommandType(commands.SettleWithdrawal{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (WithdrawalProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		return x.TransactionID, true, nil
	case events.FundsHeldForWithdrawal:
		return x.TransactionID, true, nil
	case events.WithdrawalDeclined:
		return x.TransactionID, true, nil
	case events.DailyDebitLimitConsumed:
		return x.TransactionID, true, nil
	case events.DailyDebitLimitExceeded:
		return x.TransactionID, true, nil
	case events.AccountDebitedForWithdrawal:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (WithdrawalProcess) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		s.Begin()

		// TODO(KM): Schedule a timeout for the requested time...
		// s.ScheduleTimeout(x, x.RequestedTransactionTimestamp)

		s.ExecuteCommand(commands.HoldFundsForWithdrawal{
			TransactionID:                 x.TransactionID,
			AccountID:                     x.AccountID,
			Amount:                        x.Amount,
			RequestedTransactionTimestamp: x.RequestedTransactionTimestamp,
		})

	case events.FundsHeldForWithdrawal:
		s.ExecuteCommand(commands.ConsumeDailyDebitLimit{
			TransactionID:                 x.TransactionID,
			AccountID:                     x.AccountID,
			Amount:                        x.Amount,
			RequestedTransactionTimestamp: x.RequestedTransactionTimestamp,
		})

	case events.DailyDebitLimitConsumed:
		s.ExecuteCommand(commands.SettleWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case events.DailyDebitLimitExceeded:
		s.ExecuteCommand(commands.DeclineWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
			Reason:        messages.ReasonDailyDebitLimitExceeded,
		})

	case events.AccountDebitedForWithdrawal,
		events.WithdrawalDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}