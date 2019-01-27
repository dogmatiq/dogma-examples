package main

import (
	"context"
	"os"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogmatest"
	"github.com/dogmatiq/dogmatest/engine"
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages"
)

func main() {
	app := &example.App{}
	en := dogmatest.NewEngine(app)

	messages := []dogma.Message{
		messages.OpenAccount{
			AccountID: "acct1",
			Name:      "Anna Smith",
		},
		messages.OpenAccount{
			AccountID: "acct2",
			Name:      "Bob Jones",
		},
		messages.Deposit{
			TransactionID: "txn1",
			AccountID:     "acct1",
			Amount:        10000,
		},
		messages.Transfer{
			TransactionID: "txn2",
			FromAccountID: "acct1",
			ToAccountID:   "acct2",
			Amount:        2500,
		},
	}

	for _, m := range messages {
		err := en.Dispatch(
			context.Background(),
			m,
			// engine.WithObserver(
			// 	fact.ObserverFunc(func(f fact.Fact) {
			// 		dapper.Print(f)
			// 		fmt.Print("\n\n")
			// 	}),
			// ),
			engine.EnableHandlerType(handler.ProjectionType, true),
		)
		if err != nil {
			panic(err)
		}
	}

	if err := app.GenerateAccountCSV(os.Stdout); err != nil {
		panic(err)
	}
}
