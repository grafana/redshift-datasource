import * as runtime from '@grafana/runtime';
import { mockDatasource, mockQuery } from './__mocks__/datasource';

describe('DataSource', () => {
  describe('applyTemplateVariables', () => {
    const replace = jest.fn();
    beforeEach(() => {
      jest.spyOn(runtime, 'getTemplateSrv').mockImplementation(() => ({ getVariables: jest.fn(), replace }));
    });
    it('applyTemplateVariables - query', () => {
      mockDatasource.applyTemplateVariables({ ...mockQuery, rawSQL: 'select * from bar' }, {});
      expect(replace).toBeCalledTimes(1);
      expect(replace).nthCalledWith(1, 'select * from bar', {}, 'singlequote');
    });
  });
});
