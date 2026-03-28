export const isDevelopment =
process.env.NODE_ENV !== 'production' ||
process.env.IS_DEV_BUILD === 'true';
