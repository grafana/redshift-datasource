import { test, expect } from '@grafana/plugin-e2e';

test.use({
  featureToggles: {
    alertingQueryAndExpressionsStepMode: false,
  },
});

test('should successfully create an alert rule', async ({
  alertRuleEditPage,
  page,
  readProvisionedDataSource,
  selectors,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-redshift-e2e.yaml', name: 'AWS Redshift E2E' });
  const queryA = alertRuleEditPage.getAlertRuleQueryRow('A');
  await queryA.datasource.set(ds.name);
  await page.waitForFunction(() => window.monaco);
  await queryA.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
  await page.keyboard.insertText('SELECT environment, temperature FROM public.long_format_example limit 2');
  await expect(alertRuleEditPage.evaluate()).toBeOK();
});
