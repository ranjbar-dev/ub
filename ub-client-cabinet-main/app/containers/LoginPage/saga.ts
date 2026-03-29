import { takeLatest, put, call, all } from 'redux-saga/effects';
import ActionTypes, { LoginData } from './constants';
import { isLoggingInAction } from './actions';
import {
  loginAPI,
  forgotPasswordAPI,
  getUserDataAPI,
} from 'services/security_service';
import {
  LandingPageAddress,
  LocalStorageKeys,
  StandardResponse,
} from 'services/constants';
import { loggedInAction } from 'containers/App/actions';
import { replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { getCurrenciesAPI } from 'services/address_management_service';
import { getCountriesAPI } from 'services/user_acount_service';
import {
  MessageService,
  MessageNames,
  EventMessageService,
  EventMessageNames,
} from 'services/message_service';
import { toast } from 'components/Customized/react-toastify';
import { ToastMessages } from 'services/toastService';
import { cookieConfig, CookieKeys, cookies } from 'services/cookie';
import { censor } from 'utils/formatters';
import { getFavPairAPI, getPairsListAPI } from 'services/pairs_service';
import { storage } from 'utils/storage';
import { Currency } from 'containers/AddressManagementPage/types';

export function* login(action: { type: string; payload: LoginData }) {
  yield put(isLoggingInAction(true));
  MessageService.send({
    name: MessageNames.SETLOADING,
    payload: true,
  });
  try {
    const response = yield call(loginAPI, {
      username: action.payload.username.toLowerCase().replace(/ /g, ''),
      password: action.payload.password,
      recaptcha: action.payload.recaptcha ?? 'recaptcha',
      ...(action.payload['2fa_code'] && {
        '2fa_code': action.payload['2fa_code'],
      }),
    });

    if (response.token && response.token.length > 0) {
      //@ts-ignore
      cookies.set(CookieKeys.Token, response.token, cookieConfig());
      cookies.set(
        CookieKeys.RefreshToken,
        response.refreshToken,
        cookieConfig(),
      );

      const email = censor({ value: action.payload.username, from: 1, to: 5 });
      //@ts-ignore
      cookies.set(CookieKeys.Email, email, cookieConfig());
      if (cookies.get(CookieKeys.FromLanding) === 'fromLanding') {
        cookies.remove(CookieKeys.FromLanding, {
          path: cookieConfig().path,
          domain: cookieConfig().domain,
        });
        location.replace(LandingPageAddress);
        return;
      }
      /////get currencies and countries
      const [
        countriesResponse,
        currenciesResponse,
        userDataResponse,
        favPairsResponse,
        pairs
      ]: [
          StandardResponse,
          StandardResponse,
          StandardResponse,
          StandardResponse,
          StandardResponse,
        ] = yield all([
          call(getCountriesAPI),
          call(getCurrenciesAPI),
          call(getUserDataAPI),
          call(getFavPairAPI),
          call(getPairsListAPI)
        ]);
      if (userDataResponse.data && userDataResponse.data.channelName) {
        storage.write(
          LocalStorageKeys.CHANNEL,
          userDataResponse.data.channelName,
        );
      }

      if (currenciesResponse.status === true) {
        storage.write(
          LocalStorageKeys.CURRENCIES,
          currenciesResponse.data.currencies,
        );
        const tmp = {};
        const currencies = storage.read(LocalStorageKeys.CURRENCIES);
        currencies.forEach((item: Currency) => {
          tmp[item.code] = item;
        });
        storage.write(LocalStorageKeys.CURRENCY_MAP, tmp);
      }

      if (countriesResponse.status === true) {
        storage.write(LocalStorageKeys.COUNTRIES, countriesResponse.data);
      }
      if (favPairsResponse.status === true) {
        const pairNameArray = favPairsResponse.data.map(
          (item: { name: string; id: number }) => item.name,
        );
        storage.write(LocalStorageKeys.FAV_PAIRS, pairNameArray);

        EventMessageService.send({
          name: EventMessageNames.GOT_FAV_PAIRS,
          payload: pairNameArray,
        });
      }
      MessageService.send({
        name: MessageNames.SETLOADING,
        payload: false,
      });
      yield all([put(isLoggingInAction(false)), put(loggedInAction(true))]);

      if (!action.payload.fromPopup) {
        yield put(replace(AppPages.AcountPage));
      } else {
        MessageService.send({
          name: MessageNames.CLOSE_MODAL,
        });
        toast.success('Successfully Logged In ');
      }

      //show user balance summary
      localStorage[LocalStorageKeys.SHOW_TOP_INFO] = 'false';
    } else if (response.need2fa === true) {
      yield put(isLoggingInAction(false));
      MessageService.send({
        name: MessageNames.OPEN_G2FA,
        payload: {
          username: action.payload.username.replace(/ /g, ''),
          password: action.payload.password,
          message:
            'Please Enter The 6-digit Code From Your Google Authenticator App',
        },
      });
      MessageService.send({
        name: MessageNames.RESET_RECAPTCHA,
      });
      return;
    } else {
      yield put(isLoggingInAction(false));
      MessageService.send({
        name: MessageNames.RESET_RECAPTCHA,
      });
      return;
    }
    // yield put(getUserDataAction());
  } catch (err) {
    yield put(isLoggingInAction(false));
  }
  MessageService.send({
    name: MessageNames.RESET_RECAPTCHA,
  });
}
function* forgot(action: {
  type: string;
  payload: { email: string; recaptcha: string };
}) {
  action.payload.email = action.payload.email.toLowerCase().replace(/ /g, '');
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(
      forgotPasswordAPI,
      action.payload,
    );
    if (response.status === false) {
      if (response.message && response.message.length > 0) {
        toast.error(response.message);
      }
      ToastMessages(response.data);
      MessageService.send({
        name: MessageNames.RESET_RECAPTCHA,
      });
      return;
    } else if (response.status === true) {
      toast.success('please check your email inbox');
    }
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    MessageService.send({ name: MessageNames.CLOSE_MODAL });
  } catch (error) {
    toast.error('error requesting data');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  }
  MessageService.send({
    name: MessageNames.RESET_RECAPTCHA,
  });
}
export default function* loginPageSaga() {
  yield takeLatest(ActionTypes.LOGIN_ACTION, login);
  yield takeLatest(ActionTypes.FORGOT_PASSWORD_ACTION, forgot);
}
