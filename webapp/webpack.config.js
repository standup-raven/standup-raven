const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const SentryWebpackPlugin = require('@sentry/webpack-plugin');
const buildProperties = require('../build_properties.json');
const fs = require('fs');
const propertiesParser = require('properties-file');

module.exports = {
    devtool: 'inline-source-map',
    entry: ['./src/index.jsx'],
    resolve: {
        modules: [
            'src',
            'node_modules',
        ],
        extensions: [
            '*',
            '.js',
            '.jsx',
        ],
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['env', 'react', 'stage-2'],
                    },
                },
            },
            {
                test: /\.svg$/,
                use: {
                    loader: 'svg-inline-loader',
                    options: {
                        removeSVGTagAttrs: false,
                    },
                },
            },
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader'],
            },
        ],
    },
    plugins: [
        new CopyWebpackPlugin({
            patterns: [
                {from: 'src/assets/images', to: 'static/'},
            ],
        }),
    ],
    externals: {
        react: 'React',
        redux: 'Redux',
        'react-redux': 'ReactRedux',
        'prop-types': 'PropTypes',
        'react-bootstrap': 'ReactBootstrap',
    },
    output: {
        path: path.join(__dirname, '/dist'),
        publicPath: '/',
        filename: 'main.js',
    },
};

if (buildProperties.sentry.enabled) {
    generateSentryCLIConfig(buildProperties.sentry);
    module.exports.plugins.push(
        new SentryWebpackPlugin({
            include: '.',
            ignoreFile: '.sentrycliignore',
            ignore: ['node_modules', 'webpack.config.js'],
            configFile: 'sentry.properties',
        }),
    );
}

function generateSentryCLIConfig(sentrySettings) {
    const sentryCLIConfig = {
        'defaults.url': sentrySettings.server_url,
        'defaults.org': sentrySettings.org,
        'defaults.project': sentrySettings.project,
        'auth.token': sentrySettings.auth_token,
    };

    fs.writeFileSync('./sentry.properties', propertiesParser.stringify(sentryCLIConfig));
}
