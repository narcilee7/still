import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { colors, spacing, typography } from '@still/design-system';

export function ProfileScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Profile</Text>
      <Text style={styles.hint}>Your moments and resonances.</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: spacing.lg,
  },
  title: {
    fontSize: typography.title.fontSize,
    color: colors.primary,
    marginBottom: spacing.sm,
  },
  hint: {
    fontSize: typography.description.fontSize,
    color: colors.secondary,
    textAlign: 'center',
  },
});
