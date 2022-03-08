---
id: introduction
title: Introduction
slug: /
hide_title: true
---

import useBaseUrl from '@docusaurus/useBaseUrl';

export const Logo= ({ children }) =>(
  <div
    style={{
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      padding: "72px 0",
    }}>
    <img alt="Logo" src={useBaseUrl('img/logo.svg')} width="64px" height="64px" />
    <h1
      style={{
        fontSize: "3rem",
        margin: 0,
        marginLeft: "1rem",
      }}>
      YNAB
    </h1>
  </div>
)

<Logo />

The [YNAB](https://youneedabudget.com) data source lets you visualize your personal finances in [Grafana](https://grafana.com).

Grafana is an open source platform for data visualization. Grafana started out as a tool to monitor IT systems, but has over the years evolved into a platform used to visualize and monitor anything from [beehives](https://www.hiveeyes.org/) to [avocado plants](https://grafana.com/blog/2021/03/08/how-i-built-a-monitoring-system-for-my-avocado-plant-with-arduino-and-grafana-cloud/).

With the YNAB data source for Grafana, you can now monitor your personal finances to find spending patterns and understand how they change over time.

This web site aims explains how to get started with the YNAB data source and personal finance in Grafana. If you want to contribute, head over to the [project page on GitHub](https://github.com/marcusolsson/grafana-ynab-datasource).

Now, let's get you started by first [installing the YNAB data source](installation.md).

> If you like the YNAB data source, make sure to check out [my other plugins](https://marcus.se.net/projects/) as well!
