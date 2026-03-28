/*
 * NotFoundPage Messages
 *
 * This contains all the text for the NotFoundPage component.
 */
import { defineMessages } from 'react-intl';

export const scope = 'containers.NotFoundPage';

export default defineMessages({
  header: {
    id: `${scope}.header`,
    defaultMessage: 'Page not found.',
  },
  Pagenotfound: {
    id: `${scope}.Pagenotfound`,
    defaultMessage: 'ET.Pagenotfound',
  },
  GoToHome: {
    id: `${scope}.GoToHome`,
    defaultMessage: 'ET.GoToHome',
  },
});
