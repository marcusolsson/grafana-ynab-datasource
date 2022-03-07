import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';
import { InlineField, InlineFieldRow, MultiSelect, RadioButtonGroup, Select } from '@grafana/ui';
import { defaults } from 'lodash';
import React, { useEffect, useState } from 'react';
import { DataSource } from './datasource';
import { Budget, defaultQuery, YNABDataSourceOptions, YNABQuery } from './types';

type Props = QueryEditorProps<DataSource, YNABQuery, YNABDataSourceOptions>;

export const QueryEditor = (props: Props) => {
  const { onChange, onRunQuery } = props;

  const query = defaults(props.query, defaultQuery);
  const { budgetId, accountIds } = query;

  const [budgets, setBudgets] = useState<Budget[]>([]);

  // Retrieve all budgets and accounts on initial load.
  useEffect(() => {
    const refreshBudgets = async () => {
      const budgets = await props.datasource.getResource(`budgets`);
      setBudgets(budgets);
    };
    refreshBudgets();
  }, [onChange, onRunQuery, props.datasource]);

  const dashboardVariables = getTemplateSrv()
    .getVariables()
    .map((v) => ({ label: `$${v.name}`, value: `$${v.name}` }));

  // onBudgetChange updates the budget and resets the account.
  const onBudgetChange = (value?: SelectableValue<string>) => {
    onChange({ ...query, budgetId: value?.value, accountIds: [] });
    onRunQuery();
  };

  const onAccountChange = (value?: Array<SelectableValue<string>>) => {
    onChange({ ...query, accountIds: value?.map((value) => value!.value!) });
    onRunQuery();
  };

  const selectableBudgets = budgets.map<SelectableValue<string>>((budget) => ({
    label: budget.name,
    value: budget.id,
  }));

  if (dashboardVariables.length) {
    selectableBudgets.unshift(...dashboardVariables);
  }

  // Default to the first budget in the list.
  if (!query.budgetId && selectableBudgets.length) {
    onBudgetChange(selectableBudgets[0]);
  }

  const selectedBudget = budgets.find((budget) => budget.id === getTemplateSrv().replace(query.budgetId));

  const accounts = selectedBudget ? selectedBudget.accounts.filter((account) => !account.deleted) : [];

  accounts.sort((a1, a2) => {
    return a1.name.localeCompare(a2.name);
  });

  const selectableAccounts = accounts.map<SelectableValue<string>>((account) => ({
    label: account.name,
    value: account.id,
  }));

  if (dashboardVariables.length) {
    selectableAccounts.unshift(...dashboardVariables);
  }

  console.log(budgetId, accountIds);

  return (
    <>
      <InlineFieldRow>
        <InlineField label="Budget" labelWidth={14}>
          <Select width={20} value={budgetId} options={selectableBudgets} onChange={onBudgetChange} />
        </InlineField>
        <InlineField label="Account" labelWidth={14}>
          <MultiSelect
            placeholder="All accounts"
            isClearable
            value={accountIds}
            options={selectableAccounts}
            onChange={onAccountChange}
          />
        </InlineField>
        <div className="gf-form--grow">
          <div className="gf-form-label gf-form-label--grow"></div>
        </div>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Query" labelWidth={14}>
          <RadioButtonGroup
            value={query.queryType}
            options={[
              { label: 'Balance', value: 'balance' },
              { label: 'Spending', value: 'spending' },
              { label: 'Transactions', value: 'transactions' },
            ]}
            onChange={(value) => {
              onChange({ ...query, queryType: value ?? 'balance' });
              onRunQuery();
            }}
          />
        </InlineField>
        {query.queryType === 'transactions' && (
          <InlineField labelWidth={14}>
            <RadioButtonGroup
              value={query.transactionFilter}
              options={[
                { label: 'All', value: 'all' },
                { label: 'Spending', value: 'spending' },
                { label: 'Income', value: 'income' },
              ]}
              onChange={(value) => {
                onChange({ ...query, transactionFilter: value ?? 'all' });
                onRunQuery();
              }}
            />
          </InlineField>
        )}
        {query.queryType === 'spending' && (
          <InlineField labelWidth={14}>
            <RadioButtonGroup
              value={query.spendingFilter}
              options={[
                { label: 'Spending', value: 'spending' },
                { label: 'Income', value: 'income' },
              ]}
              onChange={(value) => {
                onChange({ ...query, spendingFilter: value ?? 'spending' });
                onRunQuery();
              }}
            />
          </InlineField>
        )}
        <div className="gf-form--grow">
          <div className="gf-form-label gf-form-label--grow"></div>
        </div>
      </InlineFieldRow>
      {query.queryType === 'balance' && (
        <InlineFieldRow>
          <InlineField label="Alignment period" labelWidth={14}>
            <Select
              width={20}
              value={query.period}
              options={[
                ...dashboardVariables,
                { label: 'Daily', value: 'day' },
                { label: 'Weekly', value: 'week' },
                { label: 'Monthly', value: 'month' },
              ]}
              onChange={(value) => {
                onChange({ ...query, period: value?.value ?? 'day' });
                onRunQuery();
              }}
              disabled={true}
            />
          </InlineField>
          <div className="gf-form--grow">
            <div className="gf-form-label gf-form-label--grow"></div>
          </div>
        </InlineFieldRow>
      )}

      {query.queryType === 'spending' && (
        <InlineFieldRow>
          {query.alignBy && (
            <InlineField label="Alignment period" labelWidth={14}>
              <Select
                width={20}
                value={query.period}
                options={[
                  ...dashboardVariables,
                  { label: 'Daily', value: 'day' },
                  { label: 'Weekly', value: 'week' },
                  { label: 'Monthly', value: 'month' },
                ]}
                onChange={(value) => {
                  onChange({ ...query, period: value?.value ?? 'day' });
                  onRunQuery();
                }}
                disabled={true}
              />
            </InlineField>
          )}
          <InlineField label="Align by" labelWidth={14}>
            <Select
              width={20}
              value={query.alignBy}
              options={[
                { label: 'Account', value: 'account' },
                { label: 'Payee', value: 'payee' },
                { label: 'Category', value: 'category' },
              ]}
              onChange={(value) => {
                onChange({ ...query, alignBy: value.value });
                onRunQuery();
              }}
            />
          </InlineField>
          <div className="gf-form--grow">
            <div className="gf-form-label gf-form-label--grow"></div>
          </div>
        </InlineFieldRow>
      )}
    </>
  );
};
