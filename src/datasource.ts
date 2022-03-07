import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { YNABDataSourceOptions, YNABQuery } from './types';

export class DataSource extends DataSourceWithBackend<YNABQuery, YNABDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<YNABDataSourceOptions>) {
    super(instanceSettings);
  }

  async metricFindQuery(query: string, options?: any) {
    // Retrieve DataQueryResponse based on query.
    const response = await super.getResource(`budgets`);

    // Convert query results to a MetricFindValue[]
    const values = response.map((budget: { name: string }) => ({ text: budget.name }));

    return values;
  }
}
