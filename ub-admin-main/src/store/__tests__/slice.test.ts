import { globalActions, globalReducer, initialState, selectLoggedIn } from '../slice';
import { LocalStorageKeys } from 'services/constants';
import { RootState } from 'types';

describe('globalSlice reducer', () => {
  it('should have loggedIn: false as initial state', () => {
    const state = globalReducer(undefined, { type: '@@INIT' });
    expect(state).toEqual({ loggedIn: false });
  });

  it('setIsLoggedIn(true) sets loggedIn to true', () => {
    const state = globalReducer(initialState, globalActions.setIsLoggedIn(true));
    expect(state.loggedIn).toBe(true);
  });

  it('setIsLoggedIn(false) sets loggedIn to false', () => {
    const prevState = { loggedIn: true };
    const state = globalReducer(prevState, globalActions.setIsLoggedIn(false));
    expect(state.loggedIn).toBe(false);
  });

  it('setIsLoggedIn(false) clears all LocalStorageKeys from localStorage', () => {
    const removeItemSpy = jest.spyOn(Storage.prototype, 'removeItem');
    globalReducer({ loggedIn: true }, globalActions.setIsLoggedIn(false));

    const expectedKeys = Object.values(LocalStorageKeys);
    expectedKeys.forEach((key) => {
      expect(removeItemSpy).toHaveBeenCalledWith(key);
    });
    expect(removeItemSpy).toHaveBeenCalledTimes(expectedKeys.length);

    removeItemSpy.mockRestore();
  });

  it('setIsLoggedIn(true) does NOT clear localStorage', () => {
    const removeItemSpy = jest.spyOn(Storage.prototype, 'removeItem');
    globalReducer(initialState, globalActions.setIsLoggedIn(true));
    expect(removeItemSpy).not.toHaveBeenCalled();
    removeItemSpy.mockRestore();
  });
});

describe('selectLoggedIn selector', () => {
  const buildState = (loggedIn: boolean): RootState => ({
    global: { loggedIn },
  });

  it('returns false when loggedIn is false', () => {
    expect(selectLoggedIn(buildState(false))).toBe(false);
  });

  it('returns true when loggedIn is true', () => {
    expect(selectLoggedIn(buildState(true))).toBe(true);
  });

  it('falls back to initialState when global slice is missing', () => {
    expect(selectLoggedIn({} as RootState)).toBe(false);
  });
});
