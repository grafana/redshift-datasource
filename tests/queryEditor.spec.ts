import { test, expect } from '@grafana/plugin-e2e';

test('should provide autosuggestions', async ({ page, panelEditPage, selectors }) => {
  test.slow();
  await panelEditPage.datasource.set('AWS Redshift E2E');
  await panelEditPage.timeRange.set({ from: '2008-01-01 19:00:00', to: '2008-01-02 19:00:00' });
  await panelEditPage.setVisualization('Table');

  // The following section will verify that autocompletion is behaving as expected.
  // Throughout the composition of the SQL query, the autocompletion engine will provide appropriate suggestions.
  // In this test the first few suggestions are accepted by hitting enter which will create a basic query.
  await page.waitForFunction(() => window.monaco);
  const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.press('s');
  await expect(editor.getByLabel('SELECT <column>', { exact: true })).toBeVisible();
  await page.keyboard.press('Enter');
  await expect(editor.getByLabel('*')).toBeVisible();
  await page.keyboard.press('Enter');
  await expect(editor.getByLabel('FROM')).toBeVisible();
  await page.keyboard.press('Enter');
  await expect(editor.getByLabel('catalog_history')).toBeVisible({ timeout: 15_000 });
  await page.keyboard.press('p');
  await page.keyboard.press('u');
  await expect(editor.getByLabel('public')).toBeVisible({ timeout: 15_000 });
  await page.keyboard.press('Enter');
  await expect(editor.getByLabel('average_temperature')).toBeVisible({ timeout: 15_000 });
  await page.keyboard.press('Enter');
  await expect(editor.getByLabel('GROUP BY')).toBeVisible();
  await page.keyboard.press('Enter');
  await expect(editor.getByLabel('berlin')).toBeVisible({ timeout: 15_000 });
  await page.keyboard.press('Enter');
  await expect(editor).toContainText('SELECT * FROM public.average_temperature GROUP BY berlin');
});

test('should return data when a valid query is successfully run', async ({ page, panelEditPage, selectors }) => {
  await panelEditPage.datasource.set('AWS Redshift E2E');
  await panelEditPage.timeRange.set({ from: '2008-01-01 19:00:00', to: '2008-01-02 19:00:00' });
  await panelEditPage.setVisualization('Table');

  await page.waitForFunction(() => window.monaco);
  const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.press('ControlOrMeta+A');
  await page.keyboard.insertText(
    `select saletime as time, commission as commission from sales where $__timeFilter(time)`
  );
  await expect(panelEditPage.refreshPanel()).toBeOK();
  await expect(panelEditPage.panel.fieldNames).toHaveText(['time', 'commission'], { timeout: 15_000 });
  await expect(panelEditPage.panel.data).toContainText([
    /\d{4}(-\d{2}){2} \d{2}(:\d{2}){2}/ /* matches this pattern '2008-01-01 19:12:50' */,
    /^\d*(\.\d+)?$/ /* matches integers and decimals */,
  ]);
});
