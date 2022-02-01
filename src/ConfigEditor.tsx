import React, { ChangeEvent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { YNABDataSourceOptions, OrbitSecureJsonData } from './types';

const { SecretFormField } = LegacyForms;

export const ConfigEditor = (props: DataSourcePluginOptionsEditorProps<YNABDataSourceOptions>) => {
  const { onOptionsChange, options } = props;
  const { secureJsonFields } = options;

  const secureJsonData = (options.secureJsonData || {}) as OrbitSecureJsonData;

  const onAPITokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        apiToken: event.target.value,
      },
    });
  };

  const onResetAPIToken = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiToken: '',
      },
    });
  };

  return (
    <div className="gf-form-group">
      <div className="gf-form-inline">
        <div className="gf-form">
          <SecretFormField
            placeholder=""
            isConfigured={(secureJsonFields && secureJsonFields.apiToken) as boolean}
            value={secureJsonData.apiToken || ''}
            label="API token"
            labelWidth={8}
            inputWidth={20}
            onReset={onResetAPIToken}
            onChange={onAPITokenChange}
          />
        </div>
      </div>
    </div>
  );
};
