import { e2e } from '@grafana/e2e';

import { selectors } from '../../src/selectors';

const e2eSelectors = e2e.getSelectors(selectors.components);

/**
To run these e2e tests:
- first make sure you have access to the internal grafana redshift cluster
- set up a copy of your credentials in a provisioning/datasource/aws-redshift.yaml file
- (TODO: add test credentials to provisioning repo for symlinking)

OR if you are an external grafana contributor you can create your own cluster and use the sample data provided in the 
"Getting Started with Amazon Redshift" docs:
https://docs.aws.amazon.com/redshift/latest/gsg/cm-dev-t-load-sample-data.html
*/

type RedshiftDatasourceConfig = {
  secureJsonData: {
    accessKey: string;
    secretKey: string;
  };
  jsonData: {
    clusterIdentifier: string;
    database: string;
    dbUser: string;
    defaultRegion: string;
    managedSecret: {
      arn: string;
      name: string;
    };
  };
};
type RedshiftProvision = {
  datasources: RedshiftDatasourceConfig[];
};

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, dashboard with panel',
  scenario: () => {
    e2e()
      .readProvisions(['datasources/aws-redshift.yaml'])
      .then((RedshiftProvisions: RedshiftProvision[]) => {
        const datasource = RedshiftProvisions[0].datasources[0];

        e2e.flows.addDataSource({
          checkHealth: false,
          expectedAlertMessage: 'Data source is working',
          form: () => {
            e2eSelectors.ConfigEditor.AuthenticationProvider.input().type('Access & secret key').type('{enter}');
            e2eSelectors.ConfigEditor.AccessKey.input().type(datasource.secureJsonData.accessKey);
            e2eSelectors.ConfigEditor.SecretKey.input().type(datasource.secureJsonData.secretKey);
            e2eSelectors.ConfigEditor.DefaultRegion.input()
              .click({ force: true })
              .type(datasource.jsonData.defaultRegion)
              .type('{enter}');
            e2eSelectors.ConfigEditor.ClusterID.input().click({ force: true });
            // wait for it to load
            e2eSelectors.ConfigEditor.ClusterID.testID().contains(datasource.jsonData.clusterIdentifier);
            e2eSelectors.ConfigEditor.ClusterID.input().type(datasource.jsonData.clusterIdentifier).type('{enter}');
            e2eSelectors.ConfigEditor.Database.testID().click({ force: true }).type(datasource.jsonData.database);
            e2eSelectors.ConfigEditor.DatabaseUser.testID().click({ force: true }).type(datasource.jsonData.dbUser);
          },
          type: 'Amazon Redshift',
        });

        e2e.flows.addDashboard({
          timeRange: {
            from: '2008-01-01 19:00:00',
            to: '2008-01-02 19:00:00',
          },
        });

        e2e.flows.addPanel({
          matchScreenshot: false,
          visitDashboardAtStart: false,
          queriesForm: () => {
            e2eSelectors.QueryEditor.CodeEditor.container()
              .click({ force: true })
              .type(
                `{selectall} select saletime as time, commission as commission from sales where $__timeFilter(time)`
              );

            e2eSelectors.RefreshPicker.runButton().first().click({ force: true });
          },
        });
      });
  },
});

e2e.scenario({
  describeName: 'Smoke test - managed secret',
  itName: 'Login, create data source with a managed secret',
  scenario: () => {
    e2e()
      .readProvisions(['datasources/aws-redshift.yaml'])
      .then((RedshiftProvisions: RedshiftProvision[]) => {
        const datasource = RedshiftProvisions[0].datasources[1];

        e2e.flows.addDataSource({
          checkHealth: false,
          expectedAlertMessage: 'Data source is working',
          form: () => {
            e2eSelectors.ConfigEditor.AuthenticationProvider.input().type('Access & secret key').type('{enter}');
            e2eSelectors.ConfigEditor.AccessKey.input().type(datasource.secureJsonData.accessKey);
            e2eSelectors.ConfigEditor.SecretKey.input().type(datasource.secureJsonData.secretKey);
            e2eSelectors.ConfigEditor.DefaultRegion.input()
              .click({ force: true })
              .type(datasource.jsonData.defaultRegion)
              .type('{enter}');
            e2e().get('label').contains('AWS Secrets Manager').click({ force: true });
            e2eSelectors.ConfigEditor.ManagedSecret.input().click({ force: true });
            // wait for it to load
            // e2eSelectors.ConfigEditor.ManagedSecret.testID().contains(datasource.jsonData.managedSecret.name);
            e2eSelectors.ConfigEditor.ManagedSecret.input()
              .type(datasource.jsonData.managedSecret.name)
              .type('{enter}');
            // wait for the secret to be retrieved
            e2eSelectors.ConfigEditor.ClusterIDText.testID().should(
              'have.value',
              datasource.jsonData.clusterIdentifier
            );
            e2eSelectors.ConfigEditor.Database.testID()
              .click({ force: true })
              .type(datasource.jsonData.database, { delay: 20 });
          },
          type: 'Amazon Redshift',
        });
      });
  },
});
