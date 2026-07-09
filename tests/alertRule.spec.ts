import { test, expect } from '@grafana/plugin-e2e';

test('should successfully create an alert rule', async ({
  alertRuleEditPage,
  page,
  readProvisionedDataSource,
  selectors,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-redshift-e2e.yaml', name: 'AWS Redshift E2E' });

  // In Grafana 11.6+, the alert rule editor uses a step mode UI where individual query rows
  // are only accessible after enabling advanced mode. Enable it if the switch is present.
  // We check the DOM directly rather than calling isAdvancedModeSupported(), which relies on
  // the alertingQueryAndExpressionsStepMode feature toggle that was removed in newer Grafana builds.
  const advancedModeSwitch = alertRuleEditPage.advancedModeSwitch;
  if (await advancedModeSwitch.isVisible({ timeout: 3_000 }).catch(() => false)) {
    await advancedModeSwitch.check({ force: true });
  }

  const queryA = alertRuleEditPage.getAlertRuleQueryRow('A');
  await queryA.datasource.set(ds.name);
  await page.waitForFunction(() => window.monaco);
  await queryA.getByGrafanaSelector(selectors.components.CodeEditor.container).click();

  await page.keyboard.insertText('SELECT environment, temperature FROM public.long_format_example limit 2');
  await expect(alertRuleEditPage.evaluate()).toBeOK();
});
