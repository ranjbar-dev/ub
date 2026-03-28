import { apiService } from './apiService';
import { RequestTypes, StandardResponse, webAppAddress } from './constants';

/** Simple TTL cache for static data */
const cache = new Map<string, { data: StandardResponse; expiresAt: number }>();
const CACHE_TTL = 60 * 60 * 1000; // 1 hour

function getCached(key: string): StandardResponse | null {
  const entry = cache.get(key);
  if (entry && entry.expiresAt > Date.now()) {
    return entry.data;
  }
  cache.delete(key);
  return null;
}

function setCache(key: string, data: StandardResponse): void {
  cache.set(key, { data, expiresAt: Date.now() + CACHE_TTL });
}

/** In-flight request dedup */
const pending = new Map<string, Promise<StandardResponse>>();

interface FetchOptions {
  url: string;
  isRawUrl?: boolean;
}

function fetchWithCache(opts: FetchOptions, cacheKey: string): Promise<StandardResponse> {
  const cached = getCached(cacheKey);
  if (cached) return Promise.resolve(cached);

  const inflight = pending.get(cacheKey);
  if (inflight) return inflight;

  const promise = apiService
    .fetchData({
      data: {},
      url: opts.url,
      requestType: RequestTypes.GET,
      isRawUrl: opts.isRawUrl,
    })
    .then(response => {
      setCache(cacheKey, response);
      pending.delete(cacheKey);
      return response;
    })
    .catch(error => {
      pending.delete(cacheKey);
      throw error;
    });

  pending.set(cacheKey, promise);
  return promise;
}

/** Fetches the list of countries. Cached for 1 hour. */
export const GetCountriesAPI = (): Promise<StandardResponse> =>
  fetchWithCache({ url: `${webAppAddress}main-data/country-list`, isRawUrl: true }, 'countries');

/** Fetches the list of currencies. Cached for 1 hour. */
export const GetCurrenciesAPI = (): Promise<StandardResponse> =>
  fetchWithCache({ url: `${webAppAddress}currencies`, isRawUrl: true }, 'currencies');

/** Fetches the list of admin/manager users. Cached for 1 hour. */
export const GetManagersAPI = (): Promise<StandardResponse> =>
  fetchWithCache({ url: 'user/admins' }, 'managers');

/** Invalidates all cached static data. Call on relevant mutations. */
export function invalidateGlobalDataCache(): void {
  cache.clear();
}
