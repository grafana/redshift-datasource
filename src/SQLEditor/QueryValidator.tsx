// import { css } from '@emotion/css';
// import { formattedValueToString, getValueFormat } from '@grafana/data';
// import { Icon, IconButton, Spinner, useTheme2, HorizontalGroup, Tooltip } from '@grafana/ui';
// import { BigQueryAPI, ValidationResults } from 'api';
// import React, { useState, useMemo, useEffect } from 'react';
// import { useAsyncFn } from 'react-use';
// import useDebounce from 'react-use/lib/useDebounce';
// import { BigQueryQueryNG } from 'types';

// interface QueryValidatorProps {
//   apiClient: BigQueryAPI;
//   query: BigQueryQueryNG;
//   onValidate: (isValid: boolean) => void;
//   onFormatCode: () => void;
// }

// export function QueryValidator({ apiClient, query, onValidate, onFormatCode }: QueryValidatorProps) {
//   const [validationResult, setValidationResult] = useState<ValidationResults | null>();
//   const theme = useTheme2();
//   const valueFormatter = useMemo(() => getValueFormat('bytes'), []);

//   const styles = useMemo(() => {
//     return {
//       container: css`
//         border: 1px solid ${theme.colors.border.medium};
//         border-top: none;
//         padding: ${theme.spacing(0.5, 0.5, 0.5, 0.5)};
//         display: flex;
//         justify-content: space-between;
//         font-size: ${theme.typography.bodySmall.fontSize};
//       `,
//       error: css`
//         color: ${theme.colors.error.text};
//         font-size: ${theme.typography.bodySmall.fontSize};
//         font-family: ${theme.typography.fontFamilyMonospace};
//       `,
//       valid: css`
//         color: ${theme.colors.success.text};
//       `,
//       info: css`
//         color: ${theme.colors.text.secondary};
//       `,
//       hint: css`
//         color: ${theme.colors.text.disabled};
//         white-space: nowrap;
//         cursor: help;
//       `,
//     };
//   }, [theme]);

//   const [state, validateQuery] = useAsyncFn(
//     async (q: BigQueryQueryNG) => {
//       if (!q.location || q.rawSql.trim() === '') {
//         return null;
//       }

//       return await apiClient.validateQuery(q);
//     },
//     [apiClient]
//   );

//   const [,] = useDebounce(
//     async () => {
//       const result = await validateQuery(query);
//       if (result) {
//         setValidationResult(result);
//       }

//       return null;
//     },
//     1000,
//     [query, validateQuery]
//   );

//   useEffect(() => {
//     if (validationResult?.isError) {
//       onValidate(false);
//     }
//     if (validationResult?.isValid) {
//       onValidate(true);
//     }
//   }, [validationResult, onValidate]);

//   if (!state.value && !state.loading) {
//     return null;
//   }

//   const error = state.value?.error ? state.value.error.split(':').slice(2).join(':') : '';

//   return (
//     <div className={styles.container}>
//       <div>
//         {state.loading && (
//           <div className={styles.info}>
//             <Spinner inline={true} size={12} /> Validating query...
//           </div>
//         )}
//         {!state.loading && state.value && (
//           <>
//             <>
//               {state.value.isValid && state.value.statistics && (
//                 <div className={styles.valid}>
//                   <Icon name="check" /> This query will process{' '}
//                   <strong>{formattedValueToString(valueFormatter(state.value.statistics.TotalBytesProcessed))}</strong>{' '}
//                   when run.
//                 </div>
//               )}
//             </>

//             <>{state.value.isError && <div className={styles.error}>{error}</div>}</>
//           </>
//         )}
//       </div>
//       <div>
//         <HorizontalGroup spacing="sm">
//           <IconButton onClick={onFormatCode} name="brackets-curly" size="xs" tooltip="Format query" />
//           <Tooltip content="Hit CTRL/CMD+Return to run query">
//             <Icon className={styles.hint} name="keyboard" />
//           </Tooltip>
//         </HorizontalGroup>
//       </div>
//     </div>
//   );
// }
