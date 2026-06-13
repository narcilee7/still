import * as Sentry from '@sentry/react-native';

const dsn = process.env.EXPO_PUBLIC_SENTRY_DSN;

export function initSentry() {
  if (!dsn) return;

  Sentry.init({
    dsn,
    debug: __DEV__,
    enableNative: !__DEV__,
    environment: __DEV__ ? 'development' : 'production',
  });
}

export { Sentry };
