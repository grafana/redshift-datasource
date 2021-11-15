import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { TemplateSrv } from '@grafana/runtime';
import { RedshiftQuery } from 'types';

export const getSuggestions = (templateSrv: TemplateSrv, query: RedshiftQuery) => {
  const sugs: CodeEditorSuggestionItem[] = [
    {
      label: '$__timeEpoch',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__timeFilter',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__timeFrom',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__timeTo',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__timeGroup',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__unixEpochFilter',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__unixEpochGroup',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__schema',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${query.schema || 'public'}`,
    },
    {
      label: '$__table',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${query.table}`,
    },
    {
      label: '$__column',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${query.column}`,
    },
  ];

  templateSrv.getVariables().forEach((variable) => {
    const label = '$' + variable.name;
    let val = templateSrv.replace(label);
    if (val === label) {
      val = '';
    }
    sugs.push({
      label,
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Template Variable) ${val}`,
    });
  });

  return sugs;
};
