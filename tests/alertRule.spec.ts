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

  // Enable advanced mode if the step-mode UI is present (Grafana 11.5+).
  // In nightly, alertingQueryAndExpressionsStepMode was removed so step mode is always active.
  // Use scrollIntoViewIfNeeded to handle cases where the switch is below the fold.
  const advancedModeSwitch = alertRuleEditPage.advancedModeSwitch;
  try {
    await advancedModeSwitch.scrollIntoViewIfNeeded({ timeout: 10_000 });
    if (!await advancedModeSwitch.isChecked()) {
      await advancedModeSwitch.check({ force: true });
    }
  } catch {
    // advancedModeSwitch not present in this Grafana version
  }

  const queryA = alertRuleEditPage.getAlertRuleQueryRow('A');

  try {
    await queryA.datasource.set(ds.name);
    await page.waitForFunction(() => window.monaco);
    await queryA.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
  } catch {
    // TODO: Remove this fallback once @grafana/plugin-e2e picks up the Grafana 13 query-row selectors.
    const queryEditorRow = alertRuleEditPage.getByGrafanaSelector('data-testid Query editor row').filter({
      has: alertRuleEditPage.getByGrafanaSelector('data-testid Query editor row title A'),
    });

    // In nightly the datasource picker may sit outside the query row; fall back to page-wide search.
    const rowPicker = queryEditorRow.getByTestId(selectors.components.DataSourcePicker.inputV2);
    const datasourcePicker = (await rowPicker.isVisible({ timeout: 3_000 }).catch(() => false))
      ? rowPicker
      : page.getByTestId(selectors.components.DataSourcePicker.inputV2).first();

    await expect(datasourcePicker).toBeVisible({ timeout: 5_000 });
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
