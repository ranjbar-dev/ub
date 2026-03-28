import { LocalStorageKeys, SessionStorageKeys } from 'services/constants';
import { isDevelopment } from './environment';

export const storage = {
  write: (key: LocalStorageKeys, data: any) => {
    localStorage[key] = JSON.stringify(data);
  },
  read: (key: LocalStorageKeys, ifDoesntExist?: any) => {
    try {
      return JSON.parse(localStorage[key]);
    } catch (error) {
      if (isDevelopment) {
        console.log('error parsing local storage with key:' + key);
      }
        if (ifDoesntExist) {
          localStorage[key] = JSON.stringify(ifDoesntExist);
          return ifDoesntExist;
        }
        return null;

    }
  },
  sessionStorage: {
    write: (key: SessionStorageKeys, data: any) => {
      sessionStorage[key] = JSON.stringify(data);
    },
    read: (key: SessionStorageKeys) => {
      return JSON.parse(sessionStorage[key]);
    },
  },
};
