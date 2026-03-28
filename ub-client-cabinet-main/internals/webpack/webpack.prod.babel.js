// Important modules this config uses
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const WebpackPwaManifest = require('webpack-pwa-manifest');
const OfflinePlugin = require('offline-plugin');
const { HashedModuleIdsPlugin } = require('webpack');
const TerserPlugin = require('terser-webpack-plugin');
const CompressionPlugin = require('compression-webpack-plugin');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer')
    .BundleAnalyzerPlugin;
const InlineChunkHtmlPlugin = require("react-dev-utils/InlineChunkHtmlPlugin");

module.exports = require('./webpack.base.babel')({
    mode: 'production',

    // In production, we skip all hot-reloading stuff
    entry: [
        require.resolve('react-app-polyfill/ie11'),
        path.join(process.cwd(), 'app/app.tsx'),
    ],

    // Utilize long-term caching by adding content hashes (not compilation hashes) to compiled assets
    output: {
        filename: '[name].[chunkhash].js',
        chunkFilename: '[name].[chunkhash].chunk.js',
    },

    tsLoaders: [
        // Babel also have typescript transpiler. Uncomment this if you prefer and comment-out ts-loader
        // { loader: 'babel-loader' },
        {
            loader: 'ts-loader',
            options: {
                transpileOnly: true, // fork-ts-checker-webpack-plugin is used for type checking
                logLevel: 'info',
            },
        },
    ],

    optimization: {
        minimize: true,
        minimizer: [
            new TerserPlugin({
                terserOptions: {
                    warnings: false,
                    compress: {
                        comparisons: false,
                    },
                    parse: {},
                    mangle: true,
                    output: {
                        comments: false,
                        ascii_only: true,
                    },
                },
                parallel: true,
                cache: true,
                sourceMap: true,
            }),
        ],
        nodeEnv: 'production',
        sideEffects: true,
        concatenateModules: true,
        runtimeChunk: 'single',
        splitChunks: {
            chunks: 'all',
            maxInitialRequests: 10,
            minSize: 0,
            cacheGroups: {
                vendor: {
                    test: /[\\/]node_modules[\\/]/,
                    name(module) {
                        const packageName = module.context.match(
                            /[\\/]node_modules[\\/](.*?)([\\/]|$)/,
                        )[1];
                        return `npm.${packageName.replace('@', '')}`;
                    },
                },
            },
        },
    },

    plugins: [
        // Minify and optimize the index.html
        new HtmlWebpackPlugin({
            template: 'app/index.html',
            minify: {
                removeComments: true,
                collapseWhitespace: true,
                removeRedundantAttributes: true,
                useShortDoctype: true,
                removeEmptyAttributes: true,
                removeStyleLinkTypeAttributes: true,
                keepClosingSlash: true,
                minifyJS: true,
                minifyCSS: true,
                minifyURLs: true,
            },
            inject: true,
        }),
        new InlineChunkHtmlPlugin(HtmlWebpackPlugin, [/runtime-.+[.]js/]),
        // Put it in the end to capture all the HtmlWebpackPlugin's
        // assets manipulations and do leak its manipulations to HtmlWebpackPlugin
        new OfflinePlugin({
            relativePaths: false,
            publicPath: '/',
            appShell: '/',

            // No need to cache .htaccess. See http://mxs.is/googmp,
            // this is applied before any match in `caches` section
            excludes: ['.htaccess'],

            caches: {
                main: [':rest:'],

                // All chunks marked as `additional`, loaded after main section
                // and do not prevent SW to install. Change to `optional` if
                // do not want them to be preloaded at all (cached only when first loaded)
                additional: ['*.chunk.js'],
            },

            // Removes warning for about `additional` section usage
            safeToUseOptionalCaches: true,
        }),

        new CompressionPlugin({
            algorithm: 'gzip',
            test: /\.js$|\.css$|\.html$/,
            threshold: 10240,
            minRatio: 0.8,
        }),

        new WebpackPwaManifest({
            name: 'Client Cabinet',
            short_name: 'Client Cabinet',
            description: 'Client Cabinet project!',
            background_color: '#fafafa',
            theme_color: '#396DE0',
            inject: true,
            ios: true,
            icons: [{
                    src: path.resolve('app/images/icon-512x512.png'),
                    sizes: [72, 96, 128, 144, 192, 384, 512],
                },
                {
                    src: path.resolve('app/images/icon-512x512.png'),
                    sizes: [120, 152, 167, 180],
                    ios: true,
                },
            ],
        }),

        new HashedModuleIdsPlugin({
            hashFunction: 'sha256',
            hashDigest: 'hex',
            hashDigestLength: 20,
        }),
        new BundleAnalyzerPlugin({
            analyzerMode: 'disabled',
            generateStatsFile: true,
            statsOptions: { source: false },
        }),
    ],

    performance: {
        assetFilter: (assetFilename) =>
            !/(\.map$)|(^(main\.|favicon\.))/.test(assetFilename),
    },
});