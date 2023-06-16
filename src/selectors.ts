import { E2ESelectors } from '@grafana/e2e-selectors';

export const Components = {
  ConfigEditor: {
    AuthenticationProvider: {
      input: 'Authentication Provider',
    },
    SecretKey: {
      input: 'Secret Access Key',
    },
    AccessKey: {
      input: 'Access Key ID',
    },
    DefaultRegion: {
      input: 'Default Region',
    },
    UseServerless: {
      input: 'Serverless',
      testID: 'data-testid useServerless',
    },
    ManagedSecret: {
      input: 'Managed Secret',
      testID: 'data-testid managedSecret',
    },
    Workgroup: {
      input: 'Workgroup',
      testID: 'data-testid workgroup',
    },
    WorkgroupText: {
      input: 'Workgroup',
      testID: 'data-testid workgroup text',
    },
    ClusterID: {
      input: 'Cluster Identifier',
      testID: 'data-testid clusterID',
    },
    ClusterIDText: {
      input: 'Cluster Identifier',
      testID: 'data-testid clusterID text',
    },
    Database: {
      input: 'Database',
      testID: 'data-testid database',
    },
    DatabaseUser: {
      input: 'Database User',
      testID: 'data-testid dbuser',
    },
    schema: {
      input: 'Schema',
      testID: 'data-testid schema',
    },
    table: {
      input: 'Table',
      testID: 'data-testid table',
    },
    column: {
      input: 'Column',
      testID: 'data-testid column',
    },
    WithEvent: {
      input: 'Send events to Amazon EventBridge',
      testID: 'data-testid withEvent',
    },
  },
  QueryEditor: {
    CodeEditor: {
      container: 'Code editor container',
    },
    TableView: {
      input: 'toggle-table-view',
    },
  },
  RefreshPicker: {
    runButton: 'RefreshPicker run button',
  },
};

export const selectors: { components: E2ESelectors<typeof Components> } = {
  components: Components,
};
