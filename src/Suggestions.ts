import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { RedshiftQuery } from 'types';
import { getTemplateSrv } from '@grafana/runtime';
import { appendTemplateVariablesAsSuggestions } from '@grafana/aws-sdk';

export const getSuggestions = (query: RedshiftQuery) => {
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

  return appendTemplateVariablesAsSuggestions(getTemplateSrv, sugs);
};
