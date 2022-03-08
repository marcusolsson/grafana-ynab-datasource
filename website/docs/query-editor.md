---
id: query-editor
title: Query editor
---

This page explains the what each part of the query editor does, and how you can configure it.

### Balance

The **Balance** query returns running balances for the selected accounts.

The query starts from the first transaction and adds the amount of every subsequent transaction to get the running balance for the day of each transaction. If more than one transaction happened during a period, the value of the period is set to the balance resulting from the last transaction within that period.

For this query, we recommend the following panel configuration:

- **Visualization**: Time Series
- **Style**: Lines

### Spending

The **Spending** query groups the transactions by any of the following dimensions:

- Account
- Payee
- Category

If a period contains more than one transaction, the value of the period is set to the sum of all the transactions within that period.

For this query, we recommend the following panel configuration:

- **Visualization**: Time Series
- **Style**: Bars
- **Stack series**: Normal

### Transactions

The **Transactions** query returns all transactions within the dashboard interval.

For this query, we recommend the following panel configuration:

- **Visualization**: Table
