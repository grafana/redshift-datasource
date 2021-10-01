import React from 'react';
import { render, screen } from '@testing-library/react';
import { ResourceSelector, QueryResourceType } from './ResourceSelector';
import { select } from 'react-select-event';
import { defaultKey } from '../types';
import { selectors } from '../selectors';

const props = {
  resource: 'ManagedSecret' as QueryResourceType,
  value: null,
  list: [],
  fetch: jest.fn(),
  onChange: jest.fn(),
};

describe('AthenaResourceSelector', () => {
  it('should include a default option', () => {
    render(<ResourceSelector {...props} default="foo" value={defaultKey} />);
    expect(screen.queryByText('default (foo)')).toBeInTheDocument();
  });

  it('should select a new option', async () => {
    const onChange = jest.fn();
    const fetch = jest.fn().mockResolvedValue(['foo', 'bar']);
    render(<ResourceSelector {...props} default="foo" value={defaultKey} fetch={fetch} onChange={onChange} />);
    expect(screen.queryByText('default (foo)')).toBeInTheDocument();

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor[props.resource].input);
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, 'bar', { container: document.body });
    expect(fetch).toHaveBeenCalled();
    expect(onChange).toHaveBeenCalledWith({ label: 'bar', value: 'bar' });
  });
});
