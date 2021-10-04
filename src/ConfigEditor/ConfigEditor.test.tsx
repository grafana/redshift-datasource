import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';
import { selectors } from '../selectors';
import { mockDatasourceOptions } from '../__mocks__/datasource';
import { select } from 'react-select-event';

const secret = { name: 'foo', arn: 'arn:foo' };
const clusterIdentifier = 'cluster';
const dbUser = 'username';

jest.mock('@grafana/aws-sdk', () => {
  return {
    ...(jest.requireActual('@grafana/aws-sdk') as any),
    ConnectionConfig: function ConnectionConfig() {
      return <></>;
    },
  };
});

jest.mock('@grafana/runtime', () => {
  return {
    ...(jest.requireActual('@grafana/runtime') as any),
    getBackendSrv: () => ({
      put: jest.fn().mockResolvedValue({ datasource: {} }),
      post: jest.fn().mockResolvedValue({ dbClusterIdentifier: clusterIdentifier, username: dbUser }),
      get: jest.fn().mockResolvedValue([secret]),
    }),
  };
});

const props = mockDatasourceOptions;

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

  it('should select a secret', async () => {
    const onChange = jest.fn();
    render(<ConfigEditor {...props} onOptionsChange={onChange} />);

    screen.getByText('AWS Secrets Manager').click();

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.ManagedSecret.input);
    expect(selectEl).toBeInTheDocument();
    await select(selectEl, secret.arn, { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      jsonData: { ...props.options.jsonData, managedSecret: secret },
    });
  });

  it('should show the cluster identifier and the db user', async () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        {...props}
        onOptionsChange={onChange}
        // setting the managedSecret will trigger the secret retrieval
        options={{
          ...props.options,
          jsonData: { ...props.options.jsonData, useManagedSecret: true, managedSecret: secret },
        }}
      />
    );
    await waitFor(() => screen.getByDisplayValue(dbUser));
    // the clusterIdentifier update is delegated to the onChange function
    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      jsonData: { ...props.options.jsonData, useManagedSecret: true, managedSecret: secret, clusterIdentifier },
    });
  });
});
