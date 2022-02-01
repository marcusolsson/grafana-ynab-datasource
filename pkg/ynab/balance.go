package ynab

type Balance struct {
	Date        string
	Amount      int64
	AccountID   string
	AccountName string
}

func RunningBalance(txs []Transaction) ([]Balance, error) {
	var res []Balance

	totals := make(map[string]int64)

	for _, tx := range txs {
		if _, ok := totals[tx.AccountID]; !ok {
			totals[tx.AccountID] = 0
		}

		totals[tx.AccountID] += tx.Amount

		res = append(res, Balance{
			Date:        tx.Date,
			Amount:      totals[tx.AccountID],
			AccountID:   tx.AccountID,
			AccountName: tx.AccountName,
		})
	}

	return res, nil
}
