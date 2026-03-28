// Load the favicon and the .htaccess file
import '!file-loader?name=[name].[ext]!./images/favicon.ico';

import 'file-loader?name=.htaccess!./.htaccess';
import './index.css';

import 'components/Customized/react-toastify/dist/ReactToastify.css';

import 'react-app-polyfill/ie11';
import 'react-app-polyfill/stable';
import 'sanitize.css/sanitize.css';

// Import all the third party stuff
import * as React from 'react';
import { createRoot } from 'react-dom/client';
import { toast } from 'components/Customized/react-toastify';
import { Router } from 'react-router-dom';
import FontFaceObserver from 'fontfaceobserver';
// Import i18n messages
import { DEFAULT_LOCALE, translationMessages } from 'i18n';
import { Provider } from 'react-redux';
import history from 'utils/history';
import { createReduxHistory } from 'utils/history';

import configureStore from './configureStore';
// Import root app
import App from './containers/App/Loadable';
// Import Language Provider
import LanguageProvider from './containers/LanguageProvider';

import { OnlineStatusProvider } from 'hooks/onlineStatusHook/provider';

toast.configure({
  autoClose: 6000,
  draggable: true,
  pauseOnHover: true,
  rtl: false,
  //  pauseOnFocusLoss: false,
  // transition: Zoom,
  position: 'top-right',
  hideProgressBar: false,
});

// Observe loading of ar (to remove ar, remove the <link> tag in
// the index.html file and this observer)
const openSansObserver = new FontFaceObserver('Open Sans', {});

openSansObserver.load().then(() => {
  document.body.classList.add('fontLoaded');
  if (DEFAULT_LOCALE !== 'en') {
    if (DEFAULT_LOCALE === 'ar') {
      document.body.classList.add('arabic');
    }
  }
});

// Create redux store with history
const initialState = {};
const store = configureStore(initialState, history);
const reduxHistory = createReduxHistory(store);
const MOUNT_NODE = document.getElementById('unitedBit') as HTMLElement;
const root = createRoot(MOUNT_NODE);
const render = (messages: any, Component = App) => {
  root.render(
    <OnlineStatusProvider>
      <Provider store={store}>
        <LanguageProvider messages={messages}>
          <Router history={reduxHistory}>
            <Component />
          </Router>
        </LanguageProvider>
      </Provider>
    </OnlineStatusProvider>,
  );
};

if (module.hot) {
  module.hot.accept(['./i18n', './containers/App'], () => {
    root.unmount();
    const App = require('./containers/App').default;
    render(translationMessages, App);
  });
}

// Chunked polyfill for browsers without Intl support
if (!(window as any).Intl) {
  new Promise((resolve) => {
    resolve(import('intl'));
  })
    .then(() =>
      Promise.all([
        import('intl/locale-data/jsonp/en.js'),
        import('intl/locale-data/jsonp/de.js'),
        import('intl/locale-data/jsonp/fa.js'),
      ]),
    )
    .then(() => render(translationMessages))
    .catch((err) => {
      throw err;
    });
} else {
  render(translationMessages);
}

// ServiceWorker registration is handled by workbox-webpack-plugin in production build
//if ('serviceWorker' in navigator) {
//  navigator.serviceWorker.register('./serviceWorkers.js').then(() => {
//    console.log('service worker registered!');
//  });
//}
