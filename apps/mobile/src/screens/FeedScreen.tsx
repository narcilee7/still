import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { colors, spacing, typography } from '@still/design-system';

export function FeedScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Feed</Text>
      <Text style={styles.hint}>Moments from people you don't know.</Text>
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
