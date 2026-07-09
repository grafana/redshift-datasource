import { expect, test } from '@grafana/plugin-e2e';

test('should successfully create a variable', async ({ variableEditPage, page, selectors }) => {
  await variableEditPage.datasource.set('AWS Redshift E2E');
  await page.waitForFunction(() => window.monaco);
  const editor = page.getByTestId(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.insertText('SELECT catname FROM public.category ORDER BY catname');
  const queryDataRequest = variableEditPage.waitForQueryDataRequest();
  await variableEditPage.runQuery();
  await queryDataRequest;

  // Grafana 13+ shows variable preview as a table; older versions show a label list.
  // Use or() so the assertion works regardless of which format is rendered.
  const previewTable = variableEditPage.getByGrafanaSelector(
    selectors.pages.Dashboard.Settings.Variables.Edit.CustomVariable.previewTable
  );
  const classicalLabel = variableEditPage
    .getByGrafanaSelector(selectors.pages.Dashboard.Settings.Variables.Edit.General.previewOfValuesOption)
    .filter({ hasText: 'Classical' });

  await expect(previewTable.or(classicalLabel)).toContainText('Classical', { timeout: 15_000 });
});
