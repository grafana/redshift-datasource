import { expect, test } from '@grafana/plugin-e2e';

test('should successfully create a variable', async ({ variableEditPage, page, selectors }) => {
  await variableEditPage.datasource.set('AWS Redshift E2E');
  await page.waitForFunction(() => window.monaco);
  const editor = page.getByTestId(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.insertText('SELECT catname FROM public.category');
  const queryDataRequest = variableEditPage.waitForQueryDataRequest();
  await variableEditPage.runQuery();
  await queryDataRequest;
  await expect(variableEditPage).toDisplayPreviews([
    'MLB',
    'NFL',
    'Musicals',
    'Opera',
    'Classical',
    'NHL',
    'NBA',
    'MLS',
    'Plays',
    'Pop',
    'Jazz',
  ]);
});
