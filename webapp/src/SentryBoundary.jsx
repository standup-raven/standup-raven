import * as Sentry from '@sentry/browser';

class SentryBoundary {
    componentDidCatch(error, errorInfo) {
        Sentry.withScope((scope) => {
            Object.keys(errorInfo).forEach((key) => {
                scope.setExtra(key, errorInfo[key]);
            });
            Sentry.captureException(error);
        });
    }
}

export default SentryBoundary;
