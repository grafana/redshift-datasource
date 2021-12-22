import { getSuggestions } from 'Suggestions';
import { mockQuery } from './__mocks__/datasource';

const templateSrv = {
  getVariables: jest.fn().mockReturnValue([{ name: 'foo' }, { name: 'bar' }]),
  replace: jest.fn(),
};

jest.mock('@grafana/runtime', () => {
  return {
    ...(jest.requireActual('@grafana/runtime') as any),
    getTemplateSrv: () => templateSrv,
  };
});

describe('getSuggestions', () => {
  const macros = [
    '$__timeEpoch',
    '$__timeFilter',
    '$__timeFrom',
    '$__timeTo',
    '$__timeGroup',
    '$__unixEpochFilter',
    '$__unixEpochGroup',
    '$__schema',
    '$__table',
    '$__column',
  ];
  it('should return the list of macros', () => {
    expect(getSuggestions(mockQuery).map((s) => s.label)).toEqual(expect.arrayContaining(macros));
  });

  it('should return the list of template variables', () => {
    expect(getSuggestions(mockQuery).map((s) => s.label)).toEqual(expect.arrayContaining(['$foo', '$bar']));
  });
});
