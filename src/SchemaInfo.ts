import { TemplateSrv } from '@grafana/runtime';
import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';

export class SchemaInfo {
  constructor(private templateSrv?: TemplateSrv) {}

  getSuggestions = (): CodeEditorSuggestionItem[] => {
    const sugs: CodeEditorSuggestionItem[] = [
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
    ];

    if (this.templateSrv) {
      this.templateSrv.getVariables().forEach((variable) => {
        const label = '$' + variable.name;
        let val = this.templateSrv!.replace(label);
        if (val === label) {
          val = '';
        }
        sugs.push({
          label,
          kind: CodeEditorSuggestionItemKind.Text,
          detail: `(Template Variable) ${val}`,
        });
      });
    }

    return sugs;
  };
}
