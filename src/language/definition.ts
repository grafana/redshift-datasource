import { LanguageDefinition } from '@grafana/experimental';
import { conf, language } from './language';

const redshiftLanguageDefinition: LanguageDefinition & { id: string } = {
  id: 'redshift',
  // TODO: Load language using code splitting instead: loader: () => import('./language'),
  loader: () => Promise.resolve({ conf, language }),
};

export default redshiftLanguageDefinition;
