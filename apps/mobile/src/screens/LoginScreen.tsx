import React, { useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';
import { useOAuth } from '@clerk/expo';
import { SafeAreaView } from 'react-native-safe-area-context';
import { QuietButton, colors, spacing, typography } from '@still/design-system';

export function LoginScreen() {
  const { t } = useTranslation();
  const { startOAuthFlow } = useOAuth({ strategy: 'oauth_google' });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSignIn = async () => {
    setLoading(true);
    setError(null);
    try {
      const { createdSessionId, setActive } = await startOAuthFlow();
      if (createdSessionId) {
        await setActive?.({ session: createdSessionId });
      }
    } catch (err) {
      console.error('OAuth sign in failed', err);
      setError(t('auth.signInError'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <View style={styles.content}>
        <View style={styles.brand}>
          <Text style={styles.title}>{t('auth.title')}</Text>
          <Text style={styles.subtitle}>{t('auth.subtitle')}</Text>
        </View>
        <View style={styles.actions}>
          <QuietButton
            title={t('auth.continueWithGoogle')}
            onPress={handleSignIn}
            disabled={loading}
          />
          {error ? <Text style={styles.error}>{error}</Text> : null}
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    flex: 1,
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  brand: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  title: {
    fontSize: 48,
    lineHeight: 56,
    color: colors.primary,
    marginBottom: spacing.md,
  },
  subtitle: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.secondary,
    textAlign: 'center',
  },
  actions: {
    gap: spacing.md,
  },
  error: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.error,
    textAlign: 'center',
  },
});
