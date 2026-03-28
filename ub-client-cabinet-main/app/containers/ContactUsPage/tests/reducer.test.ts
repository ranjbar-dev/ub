import contactUsPageReducer from '../reducer';
// import { someAction } from '../actions';
import { ContainerState } from '../types';

describe('contactUsPageReducer', () => {
  let state: ContainerState;
  beforeEach(() => {
    state = {
      default: null,
      counterValue: 0,
      inputValue: '',
    };
  });

  it('returns the initial state', () => {
    const expectedResult = state;
    expect(contactUsPageReducer(undefined, {} as any)).toEqual(expectedResult);
  });

  /**
   * Example state change comparison
   *
   * it('should handle the someAction action correctly', () => {
   *   const expectedResult = {
   *     loading = true;
   *   );
   *
   *   expect(appReducer(state, someAction())).toEqual(expectedResult);
   * });
   */
});
