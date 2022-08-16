import {
  ColumnDefinition,
  getStandardSQLCompletionProvider,
  LanguageCompletionProvider,
  TableDefinition,
  TableIdentifier,
  // getStandardSQLCompletionProvider,
} from '@grafana/experimental';
import { MACROS } from './macros';

interface CompletionProviderGetterArgs {
  getColumns: React.MutableRefObject<(table: string, schema?: string) => Promise<ColumnDefinition[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<TableDefinition[]>>;
}

export const getRedshiftCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getColumns, getTables }) =>
  (monaco, language) => {
    return {
      // get standard SQL completion provider which will resolve functions and macros
      ...(language && getStandardSQLCompletionProvider(monaco, language)),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      tables: {
        resolve: async (t: TableIdentifier) => {
          return await getTables.current(t?.schema);
        },
      },
      columns: {
        resolve: async (t: TableIdentifier) => getColumns.current(t.table!, t.schema),
      },
      supportedMacros: () => MACROS,
    };
  };
