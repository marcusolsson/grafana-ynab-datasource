import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { YNABDataSourceOptions, YNABQuery } from './types';

export class DataSource extends DataSourceWithBackend<YNABQuery, YNABDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<YNABDataSourceOptions>) {
    super(instanceSettings);
  }
}
