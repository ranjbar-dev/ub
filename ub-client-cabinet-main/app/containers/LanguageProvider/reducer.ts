/*
 *
 * LanguageProvider reducer
 *
 */
import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';
import { DEFAULT_LOCALE } from '../../i18n';

export const initialState: ContainerState = {
  locale: DEFAULT_LOCALE,
};

function languageProviderReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.CHANGE_LOCALE:
      if (action.payload === 'fa') {
        document.body.classList.add('persian');
      } else {
        document.body.classList.remove('persian');
      }
      return {
        locale: action.payload,
      };
    default:
      return state;
  }
}
export default languageProviderReducer;
