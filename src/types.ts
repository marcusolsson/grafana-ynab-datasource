import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface YNABQuery extends DataQuery {
  budgetId?: string;
  accountId?: string;
  groupByCategory: boolean;
  period: string;
  alignBy?: string;
  queryType: string;
}

export const defaultQuery: Partial<YNABQuery> = {
  groupByCategory: false,
  period: 'day',
  queryType: 'spending',
};

export interface YNABDataSourceOptions extends DataSourceJsonData {}

export interface OrbitSecureJsonData {
  apiToken?: string;
}
