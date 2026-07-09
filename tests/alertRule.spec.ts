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

  // In nightly Grafana, alertingQueryAndExpressionsStepMode was removed and step mode is
  // always enabled. When the advanced mode switch is present despite our featureToggle override,
  // enable it so per-row datasource pickers are accessible.
  // On versions where the featureToggle works, the switch is not rendered, so this is a no-op.
  const advancedModeSwitch = alertRuleEditPage.advancedModeSwitch;
  if (await advancedModeSwitch.isVisible({ timeout: 5_000 }).catch(() => false)) {
    await advancedModeSwitch.check({ force: true });
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
