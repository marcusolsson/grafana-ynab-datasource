---
id: variables
title: Dashboard variables
---

Grafana lets you share configuration across multiple panels in a dashboard, using _dashboard variables_. A dashboard variable can be used to dynamically update every panel in your dashboard.

## Step 1: Create a dashboard variable

1. Click **Dashboard settings** (cog icon) at the top-right corner of your dashboard.
1. In the **Variables** tab, click **New** to create a new variable.
1. In **Name**, enter "budget".
1. In **Type**, select **Query**.
1. In **Data source**, select the YNAB data source.
1. In **Query**, enter "budgets()", and click outside the text box to update to preview at the bottom.
1. Click **Update** to save the variable.
1. Click the arrow in the top-left corner to go back to the dashboard.

In the top-left corner of the dashboard, you can now see a dropdown that contains all your budgets. You can select a budget from the dropdown to update the currently selected budget. However, you might notice that nothing changes in your panels. This is because the panels aren't using the variable yet.

## Step 2: Configure a panel to use a dashboard variable

1. Open the panel editor by clicking the panel title, and then selecting **Edit**.
1. In the query editor for the YNAB data source, click **Budget** and select "$budget".
1. Click **Apply** in the top-right corner of the page.

Now, when you change the value of the dashboard variable, that change is reflected in the panel as well. If you use the variable in all of your panels, you can quickly update your entire dashboard to show the data for that budget.

## Step 3: Add a dashboard variable for accounts

1. Add another dashboard variable called "account".
1. In **Query**, enter "accounts($budget)". "$budget" means that the variable will list all the accounts for the currently selected budget.
1. Enable **Multi-value** to be able to select multiple accounts.
1. Save the variable, and then open the panel editor.
1. In the query editor for the YNAB data source, click **Account** and select "$account".

## Step 4: Add a dashboard variable for periods

1. Add another dashboard variable called "period".
1. In **Query**, enter "periods()".
1. Save the variable, and then open the panel editor.
1. In the query editor for the YNAB data source, click **Period** and select "$period".

You can now change the period for all panels at the same time.
