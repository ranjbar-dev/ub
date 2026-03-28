import { isDevelopment } from './environment';

export const log = (toLog: any) => {
  if (isDevelopment) {
    console.log(toLog);
  }
};
