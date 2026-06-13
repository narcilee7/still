import React from 'react';
import { StatusBar } from 'expo-status-bar';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { ClerkProvider } from '@clerk/expo';
import { AuthGate } from './src/components/AuthGate';
import { TokenBridge } from './src/components/TokenBridge';
import { clerkTokenCache } from './src/services/clerkTokenCache';
import { initSentry, Sentry } from './src/services/sentry';

initSentry();

const clerkPublishableKey = process.env.EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY ?? '';

function App() {
  return (
    <ClerkProvider publishableKey={clerkPublishableKey} tokenCache={clerkTokenCache}>
      <SafeAreaProvider>
        <TokenBridge />
        <AuthGate />
        <StatusBar style="auto" />
      </SafeAreaProvider>
    </ClerkProvider>
  );
}

export default Sentry.wrap(App);
