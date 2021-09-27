import { E2ESelectors } from '@grafana/e2e-selectors';

export const Components = {
  ConfigEditor: {
    SecretKey: {
      input: 'Config editor secret key input',
    },
    AccessKey: {
      input: 'Config editor access key input',
    },
    ManagedSecret: {
      input: 'Managed Secret ARN',
      testID: 'data-testid managedSecret',
    },
    ClusterID: {
      input: 'Cluster Identifier',
      testID: 'data-testid clusterID',
    },
    Database: {
      input: 'Database',
      testID: 'data-testid database',
    },
    DatabaseUser: {
      input: 'Database User',
      testID: 'data-testid dbuser',
    },
  },
  QueryEditor: {
    CodeEditor: {
      container: 'Code editor container',
    },
  },
  RefreshPicker: {
    runButton: 'RefreshPicker run button',
  },
};

export const selectors: { components: E2ESelectors<typeof Components> } = {
  components: Components,
};
