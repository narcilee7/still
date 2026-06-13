import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { colors, spacing, typography } from '@still/design-system';

export function CreateScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Add a moment</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: typography.title.fontSize,
    color: colors.primary,
  },
});
