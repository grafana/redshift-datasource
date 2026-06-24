import { expect, test } from '@grafana/plugin-e2e';
import type { Locator, Request } from '@playwright/test';

const EXPECTED_VALUES = ['Classical', 'Jazz', 'MLB', 'MLS', 'Musicals', 'NBA', 'NFL', 'NHL', 'Opera', 'Plays', 'Pop'];
const VARIABLE_QUERY = 'SELECT catname FROM public.category ORDER BY catname';
const PREVIEW_TABLE_PAGE_SIZE = 8;

function requestIncludesVariableQuery(request: Request) {
  try {
    const body = request.postDataJSON() as {
      queries?: Array<{ rawSQL?: string }>;
      targets?: Array<{ rawSQL?: string }>;
    };
    const entries = body.queries ?? body.targets ?? [];

    if (entries.some((query) => query.rawSQL?.includes('catname'))) {
      return true;
    }
  } catch {
    // Fall through to string matching below.
  }

  // Fallback for any request shape changes across Grafana versions.
  return (request.postData() ?? '').includes('catname');
}

async function assertPreviewTableContainsValues(previewTable: Locator, values: readonly string[]) {
  const table = previewTable.getByRole('table');
  await expect(table).toBeVisible({ timeout: 15_000 });

  let currentPage = 1;

  for (let index = 0; index < values.length; index++) {
    const value = values[index];
    const pageNumber = Math.floor(index / PREVIEW_TABLE_PAGE_SIZE) + 1;

    if (pageNumber !== currentPage) {
      await previewTable.getByRole('button', { name: String(pageNumber), exact: true }).click();
      currentPage = pageNumber;
    }

    await expect(table).toContainText(value, { timeout: 15_000 });
  }
}

test('should successfully create a variable', async ({ variableEditPage, page, selectors }) => {
  await variableEditPage.datasource.set('AWS Redshift E2E');
  await page.waitForFunction(() => window.monaco);
  const editor = variableEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.insertText(VARIABLE_QUERY);

  // Listen before Tab/runQuery: Grafana 10.4+ auto-runs variable queries on change.
  const queryDataRequest = variableEditPage.waitForQueryDataRequest(requestIncludesVariableQuery);
  // Commit the monaco value so the variable model includes the query before running it.
  await page.keyboard.press('Tab');
  await variableEditPage.runQuery();
  await queryDataRequest;

  const previewTable = variableEditPage.getByGrafanaSelector(
    selectors.pages.Dashboard.Settings.Variables.Edit.CustomVariable.previewTable
  );
  const previewOptions = variableEditPage.getByGrafanaSelector(
    selectors.pages.Dashboard.Settings.Variables.Edit.General.previewOfValuesOption
  );

  // Grafana 13.1+ renders multi-column variable results in a paginated preview table.
  if (await previewTable.isVisible()) {
    await assertPreviewTableContainsValues(previewTable, EXPECTED_VALUES);
  } else {
    await expect(previewOptions.first()).toBeVisible({ timeout: 15_000 });
    await expect(previewOptions).toContainText(EXPECTED_VALUES, { timeout: 15_000 });
  }
});
