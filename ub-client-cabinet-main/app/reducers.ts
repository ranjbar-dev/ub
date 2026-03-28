/**
 * Combine all reducers in this file and export the combined reducers.
 */

import { combineReducers } from 'redux';
import { routerReducer } from 'utils/history';
import globalReducer from 'containers/App/reducer';
import languageProviderReducer from 'containers/LanguageProvider/reducer';

/**
 * Merges the main reducer with the router state and dynamically injected reducers
 */
export default function createReducer (injectedReducers = {}) {
  const rootReducer = combineReducers({
    global: globalReducer,
    language: languageProviderReducer,
    router: routerReducer,
    ...injectedReducers,
  });
  const appReducer = (state, action) => {
    if (action.type === 'App/LOGGED_IN_ACTION') {
      state = undefined;
    }
    return rootReducer(state, action);
  };
  return appReducer;
}
