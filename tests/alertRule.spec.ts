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

  if ((await queryA.getByGrafanaSelector(selectors.components.DataSourcePicker.inputV2).count()) > 0) {
    await queryA.datasource.set(ds.name);
    await page.waitForFunction(() => window.monaco);
    await queryA.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
  } else {
    // TODO: Remove this fallback once @grafana/plugin-e2e picks up the Grafana 13 query-row selectors.
    const queryEditorRow = alertRuleEditPage.getByGrafanaSelector('data-testid Query editor row').filter({
      has: alertRuleEditPage.getByGrafanaSelector('data-testid Query editor row title A'),
    });
    const datasourcePicker = queryEditorRow.getByTestId(selectors.components.DataSourcePicker.inputV2);

    await expect(datasourcePicker).toBeVisible();
    await datasourcePicker.fill(ds.name);
    await page.keyboard.press('ArrowDown');
    await page.keyboard.press('ArrowUp');
    await page.keyboard.press('Enter');

    await page.waitForFunction(() => window.monaco);
    await alertRuleEditPage
      .getByGrafanaSelector(selectors.components.CodeEditor.container, { root: queryEditorRow })
      .click();
  }

  await page.keyboard.insertText('SELECT environment, temperature FROM public.long_format_example limit 2');
  await expect(alertRuleEditPage.evaluate()).toBeOK();
});
