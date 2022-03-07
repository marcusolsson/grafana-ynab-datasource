import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface Account {
  id: string;
  name: string;
  deleted: boolean;
}

export interface Budget {
  id: string;
  name: string;
  accounts: Account[];
}

export interface YNABQuery extends DataQuery {
  budgetId?: string;
  accountIds?: string[];
  groupByCategory: boolean;
  period: string;
  alignBy?: string;
  queryType: string;
  transactionFilter: string;
  spendingFilter: string;
}

export const defaultQuery: Partial<YNABQuery> = {
  groupByCategory: false,
  period: 'day',
  queryType: 'net_worth',
  transactionFilter: 'all',
  spendingFilter: 'spending',
  alignBy: 'account',
};

export interface YNABDataSourceOptions extends DataSourceJsonData {}

export interface OrbitSecureJsonData {
  accessToken?: string;
}
