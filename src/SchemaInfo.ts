import { SelectableValue } from '@grafana/data';
import { TemplateSrv } from '@grafana/runtime';
import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { RedshiftQuery } from 'types';
import { DataSource } from './datasource';

export class SchemaInfo {
  state: Partial<RedshiftQuery>;

  schemas?: Array<SelectableValue<string>>;
  tables?: Array<SelectableValue<string>>;
  columns?: Array<SelectableValue<string>>;

  constructor(private ds: DataSource, q: Partial<RedshiftQuery>, private templateSrv?: TemplateSrv) {
    this.state = { ...q };
    if (!q.schema) {
      // The default schema is "public"
      this.state.schema = 'public';
    }
  }

  updateState(state: Partial<RedshiftQuery>): Partial<RedshiftQuery> {
    // Clean up related state
    if (state.schema) {
      this.tables = undefined;
      this.columns = undefined;
    }
    if (state.table) {
      this.columns = undefined;
    }

    const merged = { ...this.state, ...state };
    if (this.templateSrv) {
      if (merged.schema) {
        merged.schema = this.templateSrv.replace(merged.schema);
      }
      if (merged.table) {
        merged.table = this.templateSrv.replace(merged.table);
      }
      if (merged.column) {
        merged.column = this.templateSrv.replace(merged.column);
      }
    }
    return (this.state = merged);
  }

  getSchemas = async (query?: string) => {
    if (this.schemas) {
      return Promise.resolve(this.schemas);
    }
    return this.ds.getResource('schemas').then((vals: string[]) => {
      this.schemas = vals.map((name) => {
        return { label: name, value: name };
      });
      this.schemas.push({
        label: '-- remove --',
        value: '',
      });
      return this.schemas;
    });
  };

  getTables = async (query?: string) => {
    if (this.tables) {
      return Promise.resolve(this.tables);
    }
    return this.ds.postResource('tables', { schema: this.state.schema || '' }).then((vals: string[]) => {
      this.tables = vals.map((name) => {
        return { label: name, value: name };
      });
      this.tables.push({
        label: '-- remove --',
        value: '',
      });
      return this.tables;
    });
  };

  getColumns = async (query?: string) => {
    if (this.columns) {
      return Promise.resolve(this.columns);
    }
    if (!this.state.table) {
      return Promise.resolve([{ label: 'table not configured', value: '' }]);
    }
    return this.ds.postResource('columns', { table: this.state.table }).then((vals: string[]) => {
      this.columns = vals.map((name) => {
        return { label: name, value: name };
      });
      this.columns.push({
        label: '-- remove --',
        value: '',
      });
      return this.columns;
    });
  };

  async preload() {
    await this.getSchemas();
    await this.getTables();
    if (this.state.table) {
      this.getColumns();
    }
  }

  getSuggestions = (): CodeEditorSuggestionItem[] => {
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
        detail: `(Macro) ${this.state.schema}`,
      },
      {
        label: '$__table',
        kind: CodeEditorSuggestionItemKind.Text,
        detail: `(Macro) ${this.state.table}`,
      },
      {
        label: '$__column',
        kind: CodeEditorSuggestionItemKind.Text,
        detail: `(Macro) ${this.state.column}`,
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
