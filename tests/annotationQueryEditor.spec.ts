import { expect, test } from '@grafana/plugin-e2e';

test('should successfully create an annotation', async ({ annotationEditPage, page, selectors }) => {
  await annotationEditPage.datasource.set('AWS Redshift E2E');
  const editor = page.getByTestId(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.insertText('SELECT * FROM public.events');
  await expect(annotationEditPage).toHaveAlert('success', { hasText: '8 events (from 4 fields)' });
});
