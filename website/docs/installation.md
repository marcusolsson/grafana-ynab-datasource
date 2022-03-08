---
id: installation
title: Installation
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

To use the YNAB data source, you first need to have a running Grafana installation. This page explains how to install Grafana and the YNAB data source.

## Grafana Cloud

The easiest way to get started with YNAB for Grafana is using [Grafana Cloud](https://grafana.com/products/cloud). Grafana Cloud lets you run Grafana for free (with some limitations).

### Step 1: Set up a Grafana Cloud instance

1. [Sign up for Grafana Cloud](https://grafana.com/auth/sign-up).
1. Browse to **My Account** -> **Overview**.
1. Scroll down to the boxes for Grafana, Prometheus, Loki, and so on.
1. In the Grafana box, click **Log In** to start Grafana.

### Step 2: Install the YNAB data source

1. Browse to the [YNAB data source](https://grafana.com/grafana/plugins/marcusolsson-ynab-datasource/) plugin page.
1. Click the **Installation** tab.
1. Select the account you want to install the plugin on.
1. Click **Install plugin** for the instance you want to install the plugin on.

It might take a few minutes before the plugin has been installed. Next, you'll [configure the YNAB data source](configuration.md).

## Self-managed Grafana

If you prefer to run Grafana yourself, follow the instructions on how to [Download](https://grafana.com/grafana/download?edition=oss).

You can install the plugin using [grafana-cli](https://grafana.com/docs/grafana/latest/administration/cli/), or by downloading the plugin manually.

### Install using grafana-cli

To install the latest version of the plugin, run the following command on the Grafana server:

<Tabs
  groupId="operating-systems"
  defaultValue="linux"
  values={[
    {label: 'Linux', value: 'linux'},
    {label: 'macOS', value: 'macos'},
    {label: 'Windows', value: 'windows'},
  ]}>
  <TabItem value="linux">

```
grafana-cli plugins install marcusolsson-ynab-datasource
```

  </TabItem>
  <TabItem value="macos">

```
grafana-cli plugins install marcusolsson-ynab-datasource
```

  </TabItem>
  <TabItem value="windows">

```
grafana-cli.exe plugins install marcusolsson-ynab-datasource
```

  </TabItem>
</Tabs>

### Install manually

1. Go to [Releases](https://github.com/marcusolsson/grafana-ynab-datasource/releases) on the GitHub project page
1. Find the release you want to install
1. Download the release by clicking the release asset called `marcusolsson-ynab-datasource-<version>.zip`. You may need to uncollapse the **Assets** section to see it.
1. Unarchive the plugin into the Grafana plugins directory

   <Tabs
     groupId="operating-systems"
     defaultValue="linux"
     values={[
       {label: 'Linux', value: 'linux'},
       {label: 'macOS', value: 'macos'},
       {label: 'Windows', value: 'windows'},
     ]}>
     <TabItem value="linux">

     ```
     unzip marcusolsson-ynab-datasource-<version>.zip
     mv marcusolsson-ynab-datasource /var/lib/grafana/plugins
     ```

     </TabItem>
     <TabItem value="macos">

     ```
     unzip marcusolsson-ynab-datasource-<version>.zip
     mv marcusolsson-ynab-datasource /usr/local/var/lib/grafana/plugins
     ```

     </TabItem>
     <TabItem value="windows">

     ```
     Expand-Archive -Path marcusolsson-ynab-datasource-<version>.zip -DestinationPath C:\grafana\data\plugins
     ```

     </TabItem>
   </Tabs>

1. Restart the Grafana server to load the plugin
