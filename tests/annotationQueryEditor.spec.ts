import { expect, test } from '@grafana/plugin-e2e';
import { gte } from 'semver';

test('should successfully create an annotation', async ({ annotationEditPage, grafanaVersion, page, selectors }) => {
  await annotationEditPage.datasource.set('AWS Redshift E2E');
  await page.waitForFunction(() => window.monaco);
  const editor = page.getByTestId(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.insertText('SELECT * FROM public.events');
  await expect(annotationEditPage.runQuery()).toBeOK();
  if (gte(grafanaVersion, '11.0.0')) {
    await expect(annotationEditPage).toHaveAlert('success', { hasText: '8 events (from 4 fields)', timeout: 15_000 });
  }
});
