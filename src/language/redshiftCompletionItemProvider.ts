import {
  ColumnDefinition,
  CompletionItemKind,
  CompletionItemPriority,
  LanguageCompletionProvider,
  LinkedToken,
  StatementPlacementProvider,
  StatementPosition,
  SuggestionKindProvider,
  TableDefinition,
  TokenType,
} from '@grafana/experimental';
import { MACROS } from './macros';

interface CompletionProviderGetterArgs {
  getColumns: React.MutableRefObject<(t: string) => Promise<ColumnDefinition[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<TableDefinition[]>>;
  getTableSchema: React.MutableRefObject<(l: string) => Promise<TableDefinition[]>>;
}

export const getRedshiftCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getColumns, getTables, getTableSchema }) =>
  () => ({
    triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
    tables: {
      resolve: async () => {
        return await getTables.current();
      },
      parseName: (token: LinkedToken) => {
        let processedToken = token;
        let tablePath = processedToken.value;

        while (processedToken.next && processedToken?.next?.type !== TokenType.Whitespace) {
          tablePath += processedToken.next.value;
          processedToken = processedToken.next;
        }
        if (tablePath.trim().startsWith('`')) {
          return tablePath.slice(1);
        }

        return tablePath;
      },
    },
    columns: {
      resolve: async (t: string) => {
        return await getColumns.current(t);
      },
    },
    customSuggestionKinds: customSuggestionKinds(getTables, getTableSchema),
    customStatementPlacement,
    supportedMacros: () => MACROS,
  });

export enum CustomStatementPlacement {
  AfterDataset = 'afterDataset',
}

export enum CustomSuggestionKind {
  SchemaOrTable = 'schemaOrTable',
}

export const customStatementPlacement: StatementPlacementProvider = () => [
  // Overriding default befaviour of AfterFrom resolver
  {
    id: StatementPosition.AfterFromKeyword,
    overrideDefault: true,
    resolve: (currentToken) => {
      const getPreviousNonWhiteSpaceToken = currentToken?.getPreviousNonWhiteSpaceToken();
      if (getPreviousNonWhiteSpaceToken?.value.toLowerCase() === 'from') {
        return true;
      } else if (currentToken?.value.toLowerCase() === '.') {
        return true;
      } else if (
        getPreviousNonWhiteSpaceToken?.value.toLowerCase() === '.' &&
        getPreviousNonWhiteSpaceToken?.getPreviousNonWhiteSpaceToken()?.value.toLowerCase() === 'from'
      ) {
        return true;
      }
      return false;
      // if
      // const untilFrom = currentToken?.getPreviousNonWhiteSpaceToken()?.is(TokenType.Keyword, 'from');
      // if (!untilFrom) {
      //   return false;
      // }
      // let q = '';
      // for (let i = untilFrom?.length - 1; i >= 0; i--) {
      //   q += untilFrom[i].value;
      // }

      // return q.startsWith('`') && q.endsWith('`');
    },
  },
];

export const customSuggestionKinds: (
  getTables: CompletionProviderGetterArgs['getTables'],
  getTableSchema: CompletionProviderGetterArgs['getTableSchema']
) => SuggestionKindProvider = (getTables, getTableSchema) => () =>
  [
    {
      id: CustomSuggestionKind.SchemaOrTable,
      applyTo: [StatementPosition.AfterFromKeyword],
      suggestionsResolver: async (ctx) => {
        if (
          ctx.currentToken?.previous?.value === '.' &&
          ctx.currentToken?.previous?.previous?.type === TokenType.Identifier
        ) {
          const t = await getTables.current(ctx.currentToken?.previous?.previous?.value);
          console.log(t);
        }
        // const tablePath = ctx.currentToken ? getTablePath(ctx.currentToken) : '';
        // const t = await getTables.current(tablePath);
        const t = await getTableSchema.current('');

        return t.map((table) => ({
          label: table.name,
          insertText: table.completion ?? table.name,
          kind: CompletionItemKind.Field,
          sortText: CompletionItemPriority.High,
          range: {
            ...ctx.range,
            startColumn: ctx.range.endColumn,
            endColumn: ctx.range.endColumn,
          },
        }));
      },
    },
  ];

// function getTablePath(token: LinkedToken) {
//   let processedToken = token;
//   let tablePath = '';
//   while (processedToken?.previous && !processedToken.previous.isWhiteSpace()) {
//     tablePath = processedToken.previous.value + tablePath;
//     processedToken = processedToken.previous;
//   }

//   if (tablePath.startsWith('`')) {
//     tablePath = tablePath.slice(1);
//   }

//   if (tablePath.endsWith('`')) {
//     tablePath = tablePath.slice(0, -1);
//   }

//   return tablePath;
// }

// function isTypingTableIn(token: LinkedToken | null, l?: boolean) {
//   if (!token) {
//     return false;
//   }
//   const tokens = token.getPreviousUntil(TokenType.Keyword, [], 'from');
//   if (!tokens) {
//     return false;
//   }

//   let path = '';
//   for (let i = tokens.length - 1; i >= 0; i--) {
//     path += tokens[i].value;
//   }

//   if (path.startsWith('`')) {
//     path = path.slice(1);
//   }

//   return path.split('.').length === 2;
// }
