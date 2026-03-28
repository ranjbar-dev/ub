module.exports = {
  preset: 'ts-jest/presets/js-with-babel',
  testEnvironment: 'jsdom',
  collectCoverageFrom: [
    'app/**/*.{js,jsx,ts,tsx}',
    '!app/**/*.test.{js,jsx,ts,tsx}',
    '!app/*/RbGenerated*/*.{js,jsx,ts,tsx}',
    '!app/app.tsx',
    '!app/global-styles.ts',
    '!app/*/*/Loadable.{js,jsx,ts,tsx}',
  ],
  coverageThreshold: {
    global: {
      statements: 98,
      branches: 91,
      functions: 98,
      lines: 98,
    },
  },
  moduleFileExtensions: ['ts', 'tsx', 'js'],
  moduleDirectories: ['node_modules', 'app'],
  moduleNameMapper: {
    '.*\\.(css|less|styl|scss|sass)$': '<rootDir>/internals/mocks/cssModule.js',
    '.*\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/internals/mocks/image.js',
    '^components/Customized/react-toastify$': '<rootDir>/internals/mocks/toastify.js',
  },
  setupFilesAfterEnv: [
    '<rootDir>/internals/testing/test-bundler.js',
    '@testing-library/jest-dom',
  ],
  testRegex: 'tests/.*\\.test\\.(js|ts(x?))$',
  testPathIgnorePatterns: ['/node_modules/', '/node_modules_fresh/', '/node_modules_new/', '/internals/templates/'],
  transform: {
    '^.+\\.(ts(x?)|js)$': ['ts-jest', { tsconfig: 'tsconfig.json' }],
  },
  snapshotSerializers: [],
  watchPlugins: [
    'jest-watch-typeahead/filename',
    'jest-watch-typeahead/testname',
  ],
};
