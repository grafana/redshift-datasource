import { appendTemplateVariablesAsSuggestions } from '@grafana/aws-sdk';
import { getTemplateSrv } from '@grafana/runtime';
import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { RedshiftQuery } from 'types';

export const getSuggestions = (query: RedshiftQuery) => {
  const sugs: CodeEditorSuggestionItem[] = [
    {
      label: '$__schema',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${query.schema || 'public'}`,
    },
  ];

  return appendTemplateVariablesAsSuggestions(getTemplateSrv, sugs);
};
