import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { select } from 'react-select-event';

import { mockDatasourceOptions } from '../__mocks__/datasource';
import { selectors } from '../selectors';
import { ConfigEditor } from './ConfigEditor';
import userEvent from '@testing-library/user-event';

const clusterIdentifier = 'cluster';
const workgroupName = 'workgroup';
const dbUser = 'username';
const provisionedSecret = { name: 'bar', arn: 'arn:bar' };
const serverlessSecret = { name: 'foo', arn: 'arn:foo' };
const provisionedSecretFetched = { dbClusterIdentifier: clusterIdentifier, username: dbUser };
const serverlessSecretFetched = { dbClusterIdentifier: '', username: dbUser };
const cluster = { clusterIdentifier, endpoint: { address: 'bar.d.e.f', port: 456 }, database: 'db2' };
const workgroup = { workgroupName, endpoint: { address: 'foo.a.b.c', port: 123 }, database: 'db1' };

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
      get: jest.fn().mockImplementation((url, args) => {
        if (url.includes('secrets')) {
          return [provisionedSecret, serverlessSecret];
        } else if (url.includes('clusters')) {
          return [cluster];
        } else if (url.includes('workgroups')) {
          return [workgroup];
        } else {
          return [];
        }
      }),
      post: jest.fn().mockImplementation((url, args) => {
        if (url.includes('secret') && args.secretARN === 'arn:bar') {
          return provisionedSecretFetched;
        } else if (url.includes('secret') && args.secretARN === 'arn:foo') {
          return serverlessSecretFetched;
        } else {
          return;
        }
      }),
    }),
  };
});

const props = mockDatasourceOptions;

describe('ConfigEditor', () => {
  it('should display Provisioned using Secrets Manager', () => {
    render(<ConfigEditor {...props} />);
    expect(screen.getByTestId(selectors.components.ConfigEditor.WorkgroupText.testID)).not.toBeVisible();
    expect(screen.queryByTestId(selectors.components.ConfigEditor.ClusterID.testID)).toBeNull();
    expect(screen.getByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeVisible();
    expect(screen.getByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeDisabled();
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.DatabaseUser.input)).toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).not.toBeDisabled();
  });

  it('should display Provisioned using Temporary credentials', async () => {
    render(<ConfigEditor {...props} />);
    await userEvent.click(screen.getByText('Temporary credentials'));
    expect(screen.getByTestId(selectors.components.ConfigEditor.WorkgroupText.testID)).not.toBeVisible();
    expect(screen.getByTestId(selectors.components.ConfigEditor.ClusterID.testID)).toBeVisible();
    expect(screen.queryByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeNull();
    expect(screen.queryByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeNull();
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).not.toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.DatabaseUser.input)).toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.DatabaseUser.input)).not.toBeDisabled();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).not.toBeDisabled();
  });

  it('should display Serverless using Secrets Manager', () => {
    render(
      <ConfigEditor
        {...{
          ...props,
          options: {
            ...props.options,
            jsonData: {
              ...props.options.jsonData,
              useServerless: true,
            },
          },
        }}
      />
    );
    expect(screen.getByTestId(selectors.components.ConfigEditor.WorkgroupText.testID)).toBeVisible();
    expect(screen.queryByTestId(selectors.components.ConfigEditor.ClusterID.testID)).toBeNull();
    expect(screen.getByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).not.toBeVisible();
    expect(screen.getByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeDisabled();
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.DatabaseUser.input)).toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).not.toBeDisabled();
  });

  it('should display Serverless using Temporary credentials', async () => {
    render(
      <ConfigEditor
        {...{
          ...props,
          options: {
            ...props.options,
            jsonData: {
              ...props.options.jsonData,
              useServerless: true,
            },
          },
        }}
      />
    );
    await userEvent.click(screen.getByText('Temporary credentials'));
    expect(screen.getByTestId(selectors.components.ConfigEditor.WorkgroupText.testID)).toBeVisible();
    expect(screen.getByTestId(selectors.components.ConfigEditor.ClusterID.testID)).not.toBeVisible();
    expect(screen.queryByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeNull();
    expect(screen.queryByTestId(selectors.components.ConfigEditor.ClusterIDText.testID)).toBeNull();
    expect(screen.getByText(selectors.components.ConfigEditor.ManagedSecret.input)).not.toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.DatabaseUser.input)).not.toBeVisible();
    expect(screen.getByText(selectors.components.ConfigEditor.Database.input)).not.toBeDisabled();
  });

  it('should select a secret', async () => {
    const onChange = jest.fn();
    render(<ConfigEditor {...props} onOptionsChange={onChange} />);

    screen.getByText('AWS Secrets Manager').click();

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.ManagedSecret.input);
    expect(selectEl).toBeInTheDocument();
    await select(selectEl, provisionedSecret.arn, { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      jsonData: { ...props.options.jsonData, managedSecret: provisionedSecret },
    });
  });

  it('should allow user to enter a database', async () => {
    const onChange = jest.fn();
    render(<ConfigEditor {...props} onOptionsChange={onChange} />);

    const dbField = screen.getByTestId('data-testid database');
    expect(dbField).toBeInTheDocument();
    fireEvent.change(dbField, { target: { value: 'abcd' } });

    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      url: '/abcd',
      jsonData: { ...props.options.jsonData, database: 'abcd' },
    });
  });

  it('should enable WithEvent when it is toggled on', async () => {
    const onChange = jest.fn();
    render(<ConfigEditor {...props} onOptionsChange={onChange} />);
    const withEventField = screen.getByTestId(selectors.components.ConfigEditor.WithEvent.testID);
    expect(withEventField).toBeInTheDocument();

    fireEvent.click(withEventField);

    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      jsonData: { ...props.options.jsonData, withEvent: true },
    });
  });

  it('should populate the `url` prop when workGroupName is selected', async () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        {...props}
        options={{
          ...props.options,
          jsonData: { ...props.options.jsonData, database: 'test-db' },
        }}
        onOptionsChange={onChange}
      />
    );

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.WorkgroupText.input);
    expect(selectEl).toBeInTheDocument();
    await select(selectEl, workgroup.workgroupName, { container: document.body });

    await waitFor(() => expect(onChange).toHaveBeenCalledTimes(1));
    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      url: 'foo.a.b.c:123/test-db',
      jsonData: { ...props.options.jsonData, database: 'test-db', workgroupName },
    });
  });

  it('should populate the `url` prop when clusterIdentifier is selected', async () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        {...props}
        options={{
          ...props.options,
          jsonData: { ...props.options.jsonData, database: 'test-db', useManagedSecret: false },
        }}
        onOptionsChange={onChange}
      />
    );

    const selectEl = screen.getByRole('combobox', { name: selectors.components.ConfigEditor.ClusterID.input });
    expect(selectEl).toBeInTheDocument();
    await select(selectEl, cluster.clusterIdentifier, { container: document.body });

    await waitFor(() => expect(onChange).toHaveBeenCalledTimes(1));
    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      url: 'bar.d.e.f:456/test-db',
      jsonData: {
        ...props.options.jsonData,
        database: 'test-db',
        clusterIdentifier: clusterIdentifier,
        useManagedSecret: false,
      },
    });
  });

  it('should update an existing url when specifying a database', async () => {
    const onChange = jest.fn();
    await act(async () => {
      render(
        <ConfigEditor
          {...props}
          onOptionsChange={onChange}
          options={{
            ...props.options,
            url: 'my.cluster.address:123/my-old-db',
            jsonData: {
              ...props.options.jsonData,
              clusterIdentifier,
              useServerless: false,
              useManagedSecret: false,
            },
          }}
        />
      );
    });

    const dbField = screen.getByTestId('data-testid database');
    expect(dbField).toBeInTheDocument();
    fireEvent.change(dbField, { target: { value: 'abcd' } });

    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith({
      ...props.options,
      // the endpoint is updated as re-fetched
      url: 'bar.d.e.f:456/abcd',
      jsonData: {
        ...props.options.jsonData,
        clusterIdentifier,
        useServerless: false,
        useManagedSecret: false,
        database: 'abcd',
      },
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
          jsonData: {
            ...props.options.jsonData,
            useServerless: false,
            useManagedSecret: true,
            managedSecret: provisionedSecret,
          },
        }}
      />
    );

    // the clusterIdentifier and dbUser update is delegated to the onChange function
    await waitFor(() =>
      expect(onChange).toHaveBeenCalledWith({
        ...props.options,
        url: 'bar.d.e.f:456/',
        jsonData: {
          ...props.options.jsonData,
          dbUser,
          useServerless: false,
          useManagedSecret: true,
          managedSecret: provisionedSecret,
          clusterIdentifier,
        },
      })
    );
  });

  it('should show the dbUser', async () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        {...props}
        onOptionsChange={onChange}
        // setting the managedSecret will trigger the secret retrieval
        options={{
          ...props.options,
          jsonData: {
            ...props.options.jsonData,
            useServerless: true,
            useManagedSecret: true,
            managedSecret: serverlessSecret,
            workgroupName,
          },
        }}
      />
    );

    // the dbUser update is delegated to the onChange function
    await waitFor(() =>
      expect(onChange).toHaveBeenCalledWith({
        ...props.options,
        url: 'foo.a.b.c:123/',
        jsonData: {
          ...props.options.jsonData,
          dbUser,
          useServerless: true,
          useManagedSecret: true,
          managedSecret: serverlessSecret,
          workgroupName,
        },
      })
    );
  });
});
