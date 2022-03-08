---
id: configuration
title: Configuration
---

This page explains how to configure and connect the YNAB data source to your YNAB account.

## Step 1: Generate a Personal Access Token

Grafana needs your permission to access your budget on your behalf.

1. Sign in to [YNAB](https://app.youneedabudget.com), and open the [Developer Settings](https://app.youneedabudget.com/settings/developer).
1. Under **Personal Access Token**, click **New Token**.
1. Enter your current password, and click **Generate**.
1. Copy the Personal Access Token at the top of the page.

## Step 2: Add a new YNAB data source

To connect to YNAB, you need to add a data source

1. In Grafana, go to **Configuration** -> **Data Sources**, and click **Add data source**.
1. In **Personal access token**, enter the token you copied from YNAB, and click **Save & test**.

## Step 3: Import a prebuilt dashboard

While you can build your own dashboard from scratch, the YNAB data source comes with a set of dashboards built by the community specifically for YNAB.

1. Click the **Dashboards** tab to list the prebuilt dashboards for the YNAB data source.
1. Click **Import** to the right of the dashboard you want to import.
1. In the side menu, hover the cursor over **Dashboards** (four squares icon), and click **Browse** to see all your configured dashboards.
1. Click one of the dashboard you imported.
