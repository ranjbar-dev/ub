import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the addressManagementPage state domain
 */

const selectAddressManagementPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by AddressManagementPage
 */

const makeSelectAddressManagementPage = () =>
  createSelector(
    selectAddressManagementPageDomain,
    substate => {
      return substate;
    },
  );
const makeSelectIsLoading = () =>
  createSelector(
    selectAddressManagementPageDomain,
    substate => {
      return substate.AddressManagementPage
        ? substate.AddressManagementPage.isLoading
        : true;
    },
  );
const makeSelectAddresses = () =>
  createSelector(
    selectAddressManagementPageDomain,
    substate => {
      return substate.AddressManagementPage
        ? substate.AddressManagementPage.withDrawArrdesses
        : [];
    },
  );
const makeSelectCurrencies = () =>
  createSelector(
    selectAddressManagementPageDomain,
    substate => {
      return substate.AddressManagementPage
        ? substate.AddressManagementPage.currencies
        : [];
    },
  );

export default makeSelectAddressManagementPage;
export {
  selectAddressManagementPageDomain,
  makeSelectIsLoading,
  makeSelectCurrencies,
  makeSelectAddresses,
};
