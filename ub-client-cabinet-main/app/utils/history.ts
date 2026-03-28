import { createBrowserHistory } from 'history';
import { createReduxHistoryContext } from 'redux-first-history';

const history = createBrowserHistory();
const { createReduxHistory, routerMiddleware, routerReducer } =
  createReduxHistoryContext({ history });

export { createReduxHistory, routerMiddleware, routerReducer };
export default history;
