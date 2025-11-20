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
  await expect(variableEditPage).toDisplayPreviews(
    ['Classical', 'Jazz', 'MLB', 'MLS', 'Musicals', 'NBA', 'NFL', 'NHL', 'Opera', 'Plays', 'Pop'],
    { timeout: 15_000 }
  );
});
