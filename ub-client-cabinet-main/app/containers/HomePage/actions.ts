import { action } from 'typesafe-actions';
// import { } from './types';

import ActionTypes from './constants';

export const defaultAction = (data: any) =>
  action(ActionTypes.DEFAULT_ACTION, data);
// export const loginAction = (payload: LoginData) =>
//   action(ActionTypes.LOGIN_ACTION, payload);
