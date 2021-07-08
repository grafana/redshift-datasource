import { e2e } from '@grafana/e2e';
import { selectors } from '../../src/selectors';

const e2eSelectors = e2e.getSelectors(selectors.components);

/**
To run these e2e tests:
- first make sure you have access to the internal grafana redshift cluster
- set up a copy of your credentials in a provisioning/datasource/redshift.yaml file
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
    clusterId: string;
    database: string;
    dbUser: string;
    defaultRegion: string;
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
      .readProvisions(['datasources/redshift.yaml'])
      .then((RedshiftProvisions: RedshiftProvision[]) => {
        const datasource = RedshiftProvisions[0].datasources[0];

        e2e.flows.addDataSource({
          checkHealth: false,
          expectedAlertMessage: 'Data source is working',
          form: () => {
            e2e()
              .get('.aws-config-authType')
              .find(`input`)
              .click({ force: true })
              .type('Access & secret key')
              .type('{enter}');
            e2eSelectors.ConfigEditor.AccessKey.input().type(datasource.secureJsonData.accessKey);
            e2eSelectors.ConfigEditor.SecretKey.input().type(datasource.secureJsonData.secretKey);
            e2e()
              .get('.aws-config-defaultRegion')
              .find(`input`)
              .click({ force: true })
              .type(datasource.jsonData.defaultRegion)
              .type('{enter}');
            e2e().get('[data-test-id="cluster-id"]').click({ force: true }).type(datasource.jsonData.clusterId);
            e2e().get('[data-test-id="database"]').click({ force: true }).type(datasource.jsonData.database);
            e2e().get('[data-test-id="dbuser"]').click({ force: true }).type(datasource.jsonData.dbUser);
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
              )
              .type('{cmd+s}');
          },
        });
      });
  },
});
