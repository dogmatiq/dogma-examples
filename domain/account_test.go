package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_Account(t *testing.T) {
	t.Run(
		"when open new account",
		func(t *testing.T) {
			t.Run(
				"it opens the account",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						ExecuteCommand(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							EventRecorded(
								events.AccountOpened{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when open existing account",
		func(t *testing.T) {
			t.Run(
				"it does not open an account that is already open",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							}).
						ExecuteCommand(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							NoneOf(
								EventTypeRecorded(events.AccountOpened{}),
							),
						)
				},
			)
		},
	)
}
