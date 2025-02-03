import { MacroType } from '@grafana/plugin-ui';

const COLUMN = 'column',
  RELATIVE_TIME_STRING = "'1m'";

export const SCHEMA_MACRO = '$__schema';
export const TABLE_MACRO = '$__table';

export const MACROS = [
  {
    id: '$__timeFilter(dateColumn)',
    name: '$__timeFilter(dateColumn)',
    text: '$__timeFilter',
    args: [COLUMN],
    type: MacroType.Filter,
    description:
      "Will be replaced by a time range filter using the specified column name. For example, time BETWEEN '2017-07-18T11:15:52Z' AND '2017-07-18T11:15:52Z'",
  },
  {
    id: '$__timeFrom()',
    name: '$__timeFrom()',
    text: '$__timeFrom',
    args: [],
    type: MacroType.Filter,
    description:
      "Will be replaced by the start of the currently active time selection. For example, '2017-07-18T11:15:52Z'",
  },
  {
    id: '$__timeTo()',
    name: '$__timeTo()',
    text: '$__timeTo',
    args: [],
    type: MacroType.Filter,
    description:
      "Will be replaced by the end of the currently active time selection. For example, '2017-07-18T11:15:52Z'",
  },
  {
    id: '$__timeEpoch(timeColumn)',
    name: '$__timeEpoch(timeColumn)',
    text: '$__timeEpoch',
    args: [COLUMN],
    type: MacroType.Filter,
    description:
      'Will be replaced by an expression to convert to a UNIX timestamp and rename the column to time. For example, UNIX_TIMESTAMP(dateColumn) as "time"',
  },
  {
    id: '$__unixEpochFilter(timeColumn)',
    name: '$__unixEpochFilter(timeColumn)',
    text: '$__unixEpochFilter',
    args: [COLUMN],
    type: MacroType.Filter,
    description:
      'Will be replaced by a time range filter using the specified column name with times represented as Unix timestamp. For example, column >= 1624406400 AND column <= 1624410000',
  },
  {
    id: "$__timeGroup(timeColumn, '1m')",
    name: "$__timeGroup(timeColumn, '1m')",
    text: '$__timeGroup',
    args: [COLUMN, RELATIVE_TIME_STRING],
    type: MacroType.Group,
    description: `Will be replace by an expression that will group timestamps so that there is only 1 point for every period on the graph. For example, 'floor(extract(epoch from time)/60)*60 AS "time"'`,
  },
  {
    id: "$__unixEpochGroup(timeColumn, '1m')",
    name: "$__unixEpochGroup(timeColumn, '1m')",
    text: '$__unixEpochGroup',
    args: [COLUMN, RELATIVE_TIME_STRING],
    type: MacroType.Group,
    description: `Will be replace by an expression that will group epoch timestamps so that there is only 1 point for every period on the graph. For example, 'floor(time/60)*60 AS "time"'`,
  },
  {
    id: '$__column',
    name: '$__column',
    text: '$__column',
    args: [],
    type: MacroType.Column,
    description: 'Will be replaced by the query column.',
  },
  {
    id: TABLE_MACRO,
    name: TABLE_MACRO,
    text: TABLE_MACRO,
    args: [],
    type: MacroType.Table,
    description: 'Will be replaced by the query table.',
  },
  {
    id: SCHEMA_MACRO,
    name: SCHEMA_MACRO,
    text: SCHEMA_MACRO,
    args: [],
    type: MacroType.Table,
    description: 'Will be replaced by the query schema.',
  },
];
