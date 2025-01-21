import { expect, test } from '@grafana/plugin-e2e';

test('should successfully create a variable', async ({ variableEditPage, page, selectors }) => {
  await variableEditPage.datasource.set('AWS Redshift E2E');
  const editor = page.getByTestId(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.insertText('SELECT catname FROM public.category');
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
