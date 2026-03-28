import CookieProvider, { CookieAttributes } from 'js-cookie';
export enum CookieKeys {
  Token = 'ubt',
  Email = 'ube',
  RefreshToken = 'rt',
  FromLanding = 'fl',
}
export const cookies = CookieProvider;
export const cookieConfig = (): CookieAttributes => {
  const date = new Date();
  date.setTime(date.getTime() + 28 * 24 * 60 * 60 * 1000);
  const isLocal = process.env.IS_LOCAL === 'true';
  return {
    path: '/',
    domain: isLocal ? 'localhost' : '.unitedbit.com',
    expires: date,
    sameSite: 'Strict',
    secure: process.env.NODE_ENV === 'production',
  };
};
