import React, { useEffect, useState } from 'react';
import { StyleSheet, View } from 'react-native';
import { useAuth } from '@clerk/expo';
import { LoadingSpinner, colors } from '@still/design-system';
import { RootNavigator } from '../navigation/RootNavigator';
import { AuthStack } from '../navigation/AuthStack';
import { useStore } from '../store/useStore';
import { getMe } from '../services/postApi';

type AuthPhase = 'loading' | 'authenticated' | 'unauthenticated';

export function AuthGate() {
  const { isLoaded, isSignedIn } = useAuth();
  const setUser = useStore((state) => state.setUser);
  const [phase, setPhase] = useState<AuthPhase>('loading');

  useEffect(() => {
    if (!isLoaded) {
      setPhase('loading');
      return;
    }

    if (!isSignedIn) {
      setPhase('unauthenticated');
      return;
    }

    let cancelled = false;
    setPhase('loading');
    getMe()
      .then((user) => {
        if (cancelled) return;
        setUser({
          ...user,
          postsCount: 0,
          resonancesCount: 0,
        });
        setPhase('authenticated');
      })
      .catch((err) => {
        if (cancelled) return;
        console.error('getMe failed', err);
        setPhase('authenticated');
      });

    return () => {
      cancelled = true;
    };
  }, [isLoaded, isSignedIn, setUser]);

  if (phase === 'loading') {
    return (
      <View style={styles.centered}>
        <LoadingSpinner size="large" />
      </View>
    );
  }

  if (phase === 'unauthenticated') {
    return <AuthStack />;
  }

  return <RootNavigator />;
}

const styles = StyleSheet.create({
  centered: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
});
