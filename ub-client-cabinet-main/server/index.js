/* eslint consistent-return:0 import/order:0 */

const express = require('express');
const helmet = require('helmet');
const morgan = require('morgan');
const logger = require('./logger');

const argv = require('./argv');
const port = require('./port');
const setup = require('./middlewares/frontendMiddleware');
const isDev = process.env.NODE_ENV !== 'production';
const ngrok =
    (isDev && process.env.ENABLE_TUNNEL) || argv.tunnel ?
    require('ngrok') :
    false;
const { resolve } = require('path');
const app = express();

// Security headers
app.use(helmet({
    contentSecurityPolicy: false, // CSP managed separately for React app
}));

// Request logging
app.use(morgan(isDev ? 'dev' : 'combined'));

// Set up a whitelist and check against it:
// var whitelist = ['http://example1.com', 'http://example2.com']
// var corsOptions = {
//   origin: function (origin, callback) {
//     if (whitelist.indexOf(origin) !== -1) {
//       callback(null, true)
//     } else {
//       callback(new Error('Not allowed by CORS'))
//     }
//   }
// }
// app.use(cors());

// If you need a backend, e.g. an API, add your custom backend-specific middleware here
// app.use('/api', myApi);

// In production we need to pass these values in instead of relying on webpack
setup(app, {
    outputPath: resolve(process.cwd(), 'build'),
    publicPath: '/',
});

// get the intended host and port number, use localhost and port 3000 if not provided
const customHost = argv.host || process.env.HOST;
const host = customHost || null; // Let http.Server use its default IPv6/4 host
const prettyHost = customHost || 'localhost';

// use the gzipped bundle
app.get('*.js', (req, res, next) => {
    req.url = req.url + '.gz'; // eslint-disable-line
    res.set('Content-Encoding', 'gzip');
    next();
});

// Global error handler
app.use((err, req, res, next) => { // eslint-disable-line no-unused-vars
    logger.error(err.stack || err.message || err);
    res.status(500).send('Internal Server Error');
});

// Start your app.
app.listen(port, host, async(err) => {
    if (err) {
        return logger.error(err.message);
    }

    // Connect to ngrok in dev mode
    if (ngrok) {
        let url;
        try {
            url = await ngrok.connect(port);
        } catch (e) {
            return logger.error(e);
        }
        logger.appStarted(port, prettyHost, url);
    } else {
        logger.appStarted(port, prettyHost);
    }
});