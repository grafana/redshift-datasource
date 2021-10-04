import React, { useState, useEffect, useMemo } from 'react';
import { SelectableValue } from '@grafana/data';
import { InlineField, Select } from '@grafana/ui';
import { selectors } from '../selectors';
import { isEqual } from 'lodash';
import { defaultKey } from 'types';

export type QueryResourceType = 'ManagedSecret';

type Props = {
  resource: QueryResourceType;
  value: string | null;
  fetch: () => Promise<Array<string | SelectableValue<string>>>;
  onChange: (e: SelectableValue<string> | null) => void;
  dependencies?: Array<string | null>;
  tooltip?: string;
  // Options only needed for QueryEditor
  default?: string;
  // Options only needed for the ConfigEditor
  title?: string;
  disabled?: boolean;
  labelWidth?: number;
  className?: string;
  saveOptions?: () => Promise<void>;
};

export function ResourceSelector(props: Props) {
  const [resource, setResource] = useState<string | null>(props.value || props.default || null);
  const [resources, setResources] = useState<Array<string | SelectableValue>>(resource ? [resource] : []);
  const [dependencies, setDependencies] = useState(props.dependencies);
  const [isLoading, setIsLoading] = useState(false);
  const [fetched, setFetched] = useState(false);
  const defaultOpts = useMemo(() => {
    const opts: Array<SelectableValue<string>> = [
      {
        label: `default (${props.default})`,
        value: defaultKey,
        description: `Default ${props.resource} set in the data source`,
      },
    ];
    if (props.value && props.value !== defaultKey) {
      opts.push({ label: props.value, value: props.value });
    }
    return opts;
  }, [props.resource, props.default, props.value]);
  const [options, setOptions] = useState<Array<SelectableValue<string>>>(props.default ? defaultOpts : []);

  useEffect(() => {
    if (resources.length) {
      const newOptions: Array<SelectableValue<string>> = props.default ? defaultOpts : [];
      resources.forEach((r) => {
        const value = typeof r === 'string' ? r : r.value;
        if (!newOptions.find((o) => o.value === value)) {
          typeof r === 'string' ? newOptions.push({ label: r, value: r }) : newOptions.push(r);
        }
      });
      setOptions(newOptions);
    } else {
      setOptions([]);
    }
  }, [resources, defaultOpts, props.default]);

  useEffect(() => {
    // A change in the dependencies cause a state clean-up
    if (!isEqual(props.dependencies, dependencies)) {
      setFetched(false);
      setResources([]);
      setResource(null);
      props.onChange(null);
      setDependencies(props.dependencies);
    }
  }, [props, dependencies]);

  const fetch = async () => {
    if (fetched) {
      return;
    }
    if (props.saveOptions) {
      await props.saveOptions();
    }
    try {
      const resources = await props.fetch();
      setResources(resources);
    } finally {
      setFetched(true);
    }
  };

  const onChange = (e: SelectableValue<string>) => {
    props.onChange(e);
    if (e.value) {
      setResource(e.value);
    }
  };
  const onClick = async () => {
    setIsLoading(true);
    try {
      await fetch();
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <InlineField
      label={selectors.components.ConfigEditor[props.resource].input}
      labelWidth={props.labelWidth}
      tooltip={props.tooltip}
    >
      <div data-testid={selectors.components.ConfigEditor[props.resource].testID} title={props.title}>
        <Select
          aria-label={selectors.components.ConfigEditor[props.resource].input}
          options={options}
          value={props.value}
          onChange={onChange}
          isLoading={isLoading}
          className={props.className || 'min-width-6'}
          disabled={props.disabled}
          onOpenMenu={onClick}
        />
      </div>
    </InlineField>
  );
}
