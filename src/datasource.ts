import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { Budget, YNABDataSourceOptions, YNABQuery } from './types';

export class DataSource extends DataSourceWithBackend<YNABQuery, YNABDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<YNABDataSourceOptions>) {
    super(instanceSettings);
  }

  async metricFindQuery(query: { rawQuery: string }, options?: any) {
    const budgetsExp = /^budgets\(\)$/; // budgets()
    const accountsExp = /^accounts\((.*)\)$/; // accounts($budget)
    const periodsExp = /^periods\(\)$/; // periods()

    const accountsArgs = query.rawQuery.match(accountsExp);

    if (budgetsExp.test(query.rawQuery)) {
      // Retrieve DataQueryResponse based on query.
      const budgets: Budget[] = await super.getResource('budgets');

      budgets.sort((a1, a2) => {
        return a1.name.localeCompare(a2.name);
      });

      // Convert query results to a MetricFindValue[]
      return budgets.map((budget) => ({ text: budget.name, value: budget.id }));
    }

    if (accountsArgs) {
      const budgets: Budget[] = await super.getResource('budgets');

      const budget = budgets.find((budget) => budget.id === getTemplateSrv().replace(accountsArgs[1]));

      if (!budget) {
        return [];
      }

      const accounts = budget.accounts.filter((account) => !account.deleted);

      accounts.sort((a1, a2) => {
        return a1.name.localeCompare(a2.name);
      });

      return accounts.map((account) => ({ text: account.name, value: account.id }));
    }

    // Not really a query, but provides a convenient function to return the supported periods as variable.
    if (periodsExp.test(query.rawQuery)) {
      return [
        {
          text: 'Daily',
          value: 'day',
        },
        {
          text: 'Weekly',
          value: 'week',
        },
        {
          text: 'Monthly',
          value: 'month',
        },
      ];
    }

    return [];
  }

  applyTemplateVariables(query: YNABQuery, scopedVars: ScopedVars): Record<string, any> {
    const apply = (text: string): string => getTemplateSrv().replace(text, scopedVars);

    const ids: string[] | undefined = query.accountIds?.flatMap((id) => {
      const values = this.getVariableValues(id, scopedVars);
      return values.length ? values : [id];
    });

    const uniqueIds = [...new Set(ids)];

    const interpolatedQuery: YNABQuery = {
      ...query,
      budgetId: apply(query.budgetId || ''),
      accountIds: uniqueIds,
      period: apply(query.period),
    };

    console.log(interpolatedQuery);

    return interpolatedQuery;
  }

  // getVariableValues returns multiple selected values as an array rather than
  // a formatted string.
  getVariableValues(variableName: string, scopedVars: ScopedVars): string[] {
    const res: string[] = [];
    getTemplateSrv().replace(variableName, scopedVars, (value: string[]) => {
      res.push(...value);
    });
    return res;
  }
}
