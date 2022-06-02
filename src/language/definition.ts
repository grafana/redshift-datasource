import { monacoTypes } from '@grafana/ui';
import { language, conf } from './language';

export type LanguageDefinition = {
  id: string;
  extensions: string[];
  aliases: string[];
  mimetypes: string[];
  loader: () => Promise<{
    language: monacoTypes.languages.IMonarchLanguage;
    conf: monacoTypes.languages.LanguageConfiguration;
  }>;
};

const redshiftLanguageDefinition: LanguageDefinition = {
  id: 'redshift',
  extensions: ['.redshift'],
  aliases: ['Redshift'],
  mimetypes: [],
  loader: () => Promise.resolve({ conf, language }),
};

export default redshiftLanguageDefinition;
