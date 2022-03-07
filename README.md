# YNAB data source for Grafana

[![CI](https://github.com/marcusolsson/grafana-ynab-datasource/actions/workflows/ci.yml/badge.svg)](https://github.com/marcusolsson/grafana-ynab-datasource/actions/workflows/ci.yml)
[![Release](https://github.com/marcusolsson/grafana-ynab-datasource/actions/workflows/release.yml/badge.svg)](https://github.com/marcusolsson/grafana-ynab-datasource/actions/workflows/release.yml)

A data source for [You Need A Budget](https://youneedabudget.com).

## Installation

The YNAB plugin isn't available on the [plugin marketplace](https://grafana.com/plugins) yet, but you can install it manually from your terminal:

```bash
grafana-cli --pluginUrl=https://github.com/marcusolsson/grafana-ynab-datasource/releases/download/v0.1.0/marcusolsson-ynab-datasource-0.1.0.zip plugins install marcusolsson-ynab-datasource
```

## Configuration

### Step 1: Generate a Personal Access Token

Grafana needs your permission to access your budget on your behalf.

1. Sign in to [YNAB](https://app.youneedabudget.com).
1. Open the [Developer Settings](https://app.youneedabudget.com/settings/developer).
1. Under **Personal Access Token**, click **New Token**.
1. Enter your current password, and click **Generate**.
1. Copy the Personal Access Token at the top of the page.

### Step 2: Add a YNAB data source

1. In Grafana, go to **Configuration** -> **Data Sources**, and click **Add data source**.
1. In **Personal access token**, enter the token you copied from YNAB, and click **Save & test**.
