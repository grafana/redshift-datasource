import { SchemaInfo } from 'SchemaInfo';

describe('SchemaInfo', () => {
  describe('getSuggestions', () => {
    const macros = ['$__timeFilter', '$__timeFrom', '$__timeTo', '$__timeGroup'];
    it('should return the list of macros', () => {
      const schema = new SchemaInfo();
      const sugs = schema.getSuggestions();
      expect(sugs.map((s) => s.label)).toEqual(macros);
    });

    it('should return the list of template variables', () => {
      const templateSrv = {
        getVariables: jest.fn().mockReturnValue([{ name: 'foo' }, { name: 'bar' }]),
        replace: jest.fn(),
      };
      const schema = new SchemaInfo(templateSrv);
      const sugs = schema.getSuggestions();
      expect(sugs.map((s) => s.label)).toEqual(macros.concat('$foo', '$bar'));
    });
  });
});
