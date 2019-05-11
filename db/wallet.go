package db

import "github.com/oxisto/titan/model"

func GetJournalEntryIDs() ([]int64, error) {
	journalIDs := []int64{}

	err := pdb.Select(&journalIDs, `SELECT id from journal ORDER BY id DESC`)

	return journalIDs, err
}

func InsertJournalEntry(entry model.JournalEntry) error {
	_, err := pdb.Exec(`INSERT INTO journal VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		entry.ID,
		entry.Amount,
		entry.Balance,
		entry.Date,
		entry.Description,
		entry.FirstPartyID,
		entry.RefType,
		entry.SecondPartyID)

	return err
}

func GetLatestTransaction() (*model.Transaction, error) {
	var transaction model.Transaction

	if err := pdb.Get(&transaction, `SELECT * FROM transactions ORDER BY "transactionID" DESC`); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func GetTransactionIDs() ([]int64, error) {
	transactionIDs := []int64{}

	err := pdb.Select(&transactionIDs, `SELECT "transactionID" from transactions ORDER BY "transactionID" DESC`)

	return transactionIDs, err
}

func InsertTransaction(transaction model.Transaction) error {
	_, err := pdb.Exec(`INSERT INTO transactions VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		transaction.TransactionID,
		transaction.ClientID,
		transaction.Date,
		transaction.IsBuy,
		transaction.JournalRefID,
		transaction.LocationID,
		transaction.Quantity,
		transaction.TypeID,
		transaction.UnitPrice)

	return err
}
