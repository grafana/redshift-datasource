import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { TemplateSrv } from '@grafana/runtime';

export const getSuggestions = ({
  templateSrv,
  schema,
  table,
  column,
}: {
  templateSrv: TemplateSrv;
  schema?: string;
  table?: string;
  column?: string;
}) => {
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
      detail: `(Macro) ${schema || 'public'}`,
    },
    {
      label: '$__table',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${table}`,
    },
    {
      label: '$__column',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${column}`,
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
