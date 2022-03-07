import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { Button, Card, InlineField, InlineFieldRow, Input, LinkButton } from '@grafana/ui';
import React, { ChangeEvent } from 'react';
import { OrbitSecureJsonData, YNABDataSourceOptions } from './types';

export const ConfigEditor = (props: DataSourcePluginOptionsEditorProps<YNABDataSourceOptions>) => {
  const { onOptionsChange, options } = props;
  const { secureJsonFields } = options;

  const secureJsonData = (options.secureJsonData || {}) as OrbitSecureJsonData;

  const onAccessTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        accessToken: event.target.value,
      },
    });
  };

  const onResetAccessToken = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        accessToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        accessToken: '',
      },
    });
  };

  const configured = !!(secureJsonFields && secureJsonFields.accessToken);

  return (
    <>
      <Card
        heading="Getting started"
        description="To connect to YNAB, you need to generate a Personal Access Token in the YNAB app."
      >
        <Card.Actions>
          <LinkButton target="_blank" variant="secondary" href="https://app.youneedabudget.com/settings/developer">
            Generate new token
          </LinkButton>
        </Card.Actions>
      </Card>
      <InlineFieldRow>
        <InlineField label="Personal access token" disabled={configured}>
          <Input
            type="password"
            value={secureJsonData.accessToken || ''}
            onChange={onAccessTokenChange}
            placeholder={configured ? 'configured' : ''}
          />
        </InlineField>
        {configured && (
          <Button onClick={onResetAccessToken} variant="secondary">
            Reset
          </Button>
        )}
      </InlineFieldRow>
    </>
  );
};
