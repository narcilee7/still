import React from 'react';
import { Pressable, StyleSheet, Text } from 'react-native';
import { colors, spacing, typography } from '../theme';

export interface ResonateButtonProps {
  count: number;
  resonated?: boolean;
  onPress?: () => void;
}

export function ResonateButton({ count, resonated, onPress }: ResonateButtonProps) {
  return (
    <Pressable
      onPress={onPress}
      style={[styles.root, resonated && styles.active]}
      accessibilityRole="button"
      accessibilityLabel={resonated ? 'You feel this' : 'I feel this'}
    >
      <Text style={[styles.label, resonated && styles.activeLabel]}>
        {resonated ? 'I feel this' : 'I feel this'}
      </Text>
      <Text style={[styles.count, resonated && styles.activeLabel]}>{count}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  root: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: colors.border,
    gap: spacing.sm,
  },
  active: {
    backgroundColor: colors.primary,
    borderColor: colors.primary,
  },
  label: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
  },
  activeLabel: {
    color: colors.white,
  },
  count: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.primary,
  },
});
