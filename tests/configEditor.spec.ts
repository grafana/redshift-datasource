import { test, expect } from '@grafana/plugin-e2e';

test('should successfully provision a data source', async ({ gotoDataSourceConfigPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-redshift-e2e.yaml', name: 'AWS Redshift E2E' });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await expect(configPage.saveAndTest()).toBeOK();
});

test('should successfully provision a data source with a managed secret', async ({
  gotoDataSourceConfigPage,
  page,
  readProvisionedDataSource,
  selectors,
}) => {
  const ds = await readProvisionedDataSource({
    fileName: 'aws-redshift-e2e.yaml',
    name: 'AWS Redshift E2E with managed secret',
  });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await page.getByLabel('AWS Secrets Manager', { exact: true }).click();
  await page.getByLabel('Managed Secret').fill('redshiftqueryeditor-cloud-datasources-redshift-cloud-datasources');
  await expect(configPage.getByGrafanaSelector(selectors.components.Select.option)).toContainText([
    'redshiftqueryeditor-cloud-datasources-redshift-cloud-datasources',
  ]);
  await page.keyboard.press('Enter');
  await expect(page.getByLabel('Cluster Identifier')).toHaveValue('redshift-cluster-grafana');
  await expect(configPage.saveAndTest()).toBeOK();
});

test('should return an error when invalid credentials are used', async ({
  createDataSourceConfigPage,
  page,
  readProvisionedDataSource,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-redshift-e2e.yaml', name: 'AWS Redshift E2E' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });
  await page.getByLabel(/^Authentication Provider/).fill('Access & secret key');
  await page.keyboard.press('Enter');
  await page.getByLabel('Access Key ID').fill('invalid_access_id');
  await page.getByLabel('Secret Access Key').fill('invalid_secret_key');
  await page.getByLabel('Default Region').fill('us-east-2');
  await page.keyboard.press('Enter');
  await expect(configPage.saveAndTest()).not.toBeOK();
});
