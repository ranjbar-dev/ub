import { AppPages } from 'app/constants';
import { replace } from 'connected-react-router';
import { takeLatest, call, put } from 'redux-saga/effects';
import { StandardResponse, LocalStorageKeys } from 'services/constants';
import {
  GetCountriesAPI,
  GetManagersAPI,
  GetCurrenciesAPI,
} from 'services/globalDataService';
import { loginAPI } from 'services/securityService';
import { GetCurrencyPairsAPI } from 'services/userManagementService';
import { globalActions } from 'store/slice';

import { actions } from './slice';

/** Strip characters that have no place in a username/password before
 *  they ever reach the network. Defence-in-depth — does NOT replace
 *  server-side validation. */
function sanitizeCredential(value: string): string {
  return value.replace(/[\x00<>"'`]/g, '').trim();
}

export function* Login(action: {
  type: string;
  payload: { username: string; password: string };
}) {
  yield put(actions.setIsLoadingAction(true));
  try {
    yield put(actions.setErrorAction(null));
    const sanitizedPayload = {
      username: sanitizeCredential(action.payload.username),
      password: sanitizeCredential(action.payload.password),
    };
    const response: StandardResponse = yield call(loginAPI, sanitizedPayload);
    if (response.token) {
      localStorage.setItem(LocalStorageKeys.ACCESS_TOKEN, response.token);
      if (response.data && (response.data as Record<string, unknown>).refresh_token) {
        localStorage.setItem(
          LocalStorageKeys.REFRESH_TOKEN,
          (response.data as Record<string, unknown>).refresh_token as string,
        );
      }

      try {
        const payloadB64 = response.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
        const jwtPayload: { role?: string } = JSON.parse(atob(payloadB64));
        yield put(globalActions.setRole(jwtPayload.role ?? null));
      } catch {
        yield put(globalActions.setRole(null));
      }

      const countriesResponse: StandardResponse = yield call(GetCountriesAPI);
      if (countriesResponse.status === true) {
        localStorage[LocalStorageKeys.COUNTRIES] = JSON.stringify(
          countriesResponse.data,
        );
      }
      const currenciesResponse: StandardResponse = yield call(GetCurrenciesAPI);
      if (currenciesResponse.status === true) {
        localStorage[LocalStorageKeys.CURRENCIES] = JSON.stringify(
          currenciesResponse.data,
        );
      }
      const managersResponse: StandardResponse = yield call(GetManagersAPI);
      if (managersResponse.status === true) {
        localStorage[LocalStorageKeys.Managers] = JSON.stringify(
          managersResponse.data,
        );
      }
      const pairsResponse: StandardResponse<{ name: string; id: number }[]> = yield call(
        GetCurrencyPairsAPI,
        {},
      );
      if (pairsResponse.status === true) {
        localStorage[LocalStorageKeys.PAIRS] = JSON.stringify(
          pairsResponse.data.map(
            (item: { name: string; id: number }) => ({
              name: item.name,
              value: item.id + '',
            }),
          ),
        );
      }
      yield put(globalActions.setIsLoggedIn(true));
      yield put(replace(AppPages.HomePage));
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Login failed. Please try again.';
    yield put(actions.setErrorAction(message));
  } finally {
    yield put(actions.setIsLoadingAction(false));
  }
}

export function* loginPageSaga() {
  yield takeLatest(actions.LoginAction.type, Login);
}