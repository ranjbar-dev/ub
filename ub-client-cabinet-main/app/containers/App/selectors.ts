/**
 * The global state selectors
 */

import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';

const selectGlobal = (state: ApplicationRootState) => {
  return state.global;
};

const selectRoute = (state: ApplicationRootState) => {
  return state.router;
};

const makeSelectLoading = () =>
  createSelector(selectGlobal, globalState => globalState.loading);

const makeSelectLoggedIn = () =>
  createSelector(selectGlobal, globalState => globalState.loggedIn);
const makeSelectTheme = () =>
  createSelector(selectGlobal, globalState => globalState.theme);

const makeSelectLocation = () =>
  createSelector(selectRoute, routeState => routeState.location);
const selectAppDomain = (state: ApplicationRootState) => {
  return state || {};
};

const makeSelectAppState = () =>
  createSelector(selectAppDomain, substate => {
    return substate;
  });
const makeSelectLanguage = () =>
  createSelector(selectAppDomain, substate => {
    return substate.language.locale;
  });

export {
  selectGlobal,
  makeSelectAppState,
  makeSelectLoading,
  makeSelectLocation,
  makeSelectLanguage,
  makeSelectLoggedIn,
  makeSelectTheme,
};
