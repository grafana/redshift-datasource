import { getSuggestions } from 'Suggestions';

const templateSrv = {
  getVariables: jest.fn().mockReturnValue([]),
  replace: jest.fn(),
};

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
    expect(getSuggestions(templateSrv).map((s) => s.label)).toEqual(macros);
  });

  it('should return the list of template variables', () => {
    const templateSrv = {
      getVariables: jest.fn().mockReturnValue([{ name: 'foo' }, { name: 'bar' }]),
      replace: jest.fn(),
    };
    expect(getSuggestions(templateSrv).map((s) => s.label)).toEqual(macros.concat('$foo', '$bar'));
  });
});
