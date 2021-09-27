import React from 'react';
import { render, screen } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';
import { DataSourceSettings } from '@grafana/data';
import { RedshiftDataSourceOptions, RedshiftDataSourceSecureJsonData } from './types';
import userEvent from '@testing-library/user-event';

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
    expect(screen.getByText('Cluster Identifier')).toBeInTheDocument();
    expect(screen.getByText('Database')).toBeInTheDocument();
    expect(screen.getByText('DB User')).toBeInTheDocument();
  });

  it('should switch to use the Secret Manager', () => {
    render(<ConfigEditor {...props} />);
    screen.getByText('AWS Secrets Manager').click();
    expect(screen.getByText('Managed Secret ARN')).toBeInTheDocument();
    expect(screen.getByText('Cluster Identifier')).toBeInTheDocument();
    expect(screen.getByText('Database')).toBeInTheDocument();
  });

  it('should clean up related state', () => {
    render(<ConfigEditor {...props} />);
    // type a user. Using a single letter since the change method is mocked so the value is not updated
    userEvent.type(screen.getByLabelText('DB User'), 'f');
    expect(props.onOptionsChange).toHaveBeenLastCalledWith({ jsonData: { dbUser: 'f' } });

    // change auth type and clean state
    screen.getByText('AWS Secrets Manager').click();
    expect(screen.getByText('Managed Secret ARN')).toBeInTheDocument();
    expect(props.onOptionsChange).toHaveBeenLastCalledWith({ jsonData: { dbUser: '' } });
  });
});
