const { TextEncoder, TextDecoder } = require('util');
global.TextEncoder = TextEncoder;
global.TextDecoder = TextDecoder;

require('react-app-polyfill/ie11');
require('react-app-polyfill/stable');

// Suppress known React 18 legacy context warnings from react-intl v2
// These will be resolved when react-intl is upgraded to v6
const originalConsoleError = console.error;
console.error = (...args) => {
  if (typeof args[0] === 'string' && args[0].includes('uses the legacy childContextTypes API')) {
    return;
  }
  originalConsoleError.apply(console, args);
};
