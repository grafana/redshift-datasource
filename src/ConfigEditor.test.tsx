import React from 'react';
import { render, screen } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';
import { DataSourceSettings } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData } from './types';
import userEvent from '@testing-library/user-event';
import { selectors } from './selectors';

jest.mock('@grafana/aws-sdk', () => {
  return {
    ...(jest.requireActual('@grafana/aws-sdk') as any),
    ConnectionConfig: function ConnectionConfig() {
      return <></>;
    },
  };
});

const props = {
  options: {
    jsonData: {},
  } as DataSourceSettings<RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData>,
  onOptionsChange: jest.fn(),
};

describe('ConfigEditor', () => {
  it('should display temporary credentials by default', () => {
    render(<ConfigEditor {...props} />);
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.DatabaseUser.input)).toBeInTheDocument();
  });

  it('should switch to use the Secret Manager', () => {
    render(<ConfigEditor {...props} />);
    screen.getByText('AWS Secrets Manager').click();
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.ClusterID.input)).toBeInTheDocument();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).toBeInTheDocument();
  });

  it('should clean up related state', () => {
    render(<ConfigEditor {...props} />);
    // type a user. Using a single letter since the change method is mocked so the value is not updated
    userEvent.type(screen.getByTestId(selectors.components.ConfigEditor.DatabaseUser.testID), 'f');
    expect(props.onOptionsChange).toHaveBeenLastCalledWith({ jsonData: { dbUser: 'f' } });

    // change auth type and clean state
    screen.getByText('AWS Secrets Manager').click();
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeInTheDocument();
    expect(props.onOptionsChange).toHaveBeenLastCalledWith({ jsonData: { dbUser: '' } });
  });
});
