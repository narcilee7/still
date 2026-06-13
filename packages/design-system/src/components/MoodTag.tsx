import React from 'react';
import { StyleSheet, Text, TextStyle } from 'react-native';
import { Mood } from '@still/shared-types';
import { moodColors, typography } from '../theme';

export interface MoodTagProps {
  mood: Mood;
  variant?: 'large' | 'small';
  style?: TextStyle;
}

export function MoodTag({ mood, variant = 'large', style }: MoodTagProps) {
  return (
    <Text
      style={[
        variant === 'large' ? styles.large : styles.small,
        { color: moodColors[mood] ?? moodColors.still },
        style,
      ]}
    >
      {mood}
    </Text>
  );
}

const styles = StyleSheet.create({
  large: {
    fontSize: typography.mood.fontSize,
    lineHeight: typography.mood.lineHeight,
    fontWeight: typography.mood.fontWeight,
    letterSpacing: typography.mood.letterSpacing,
  },
  small: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    fontWeight: typography.title.fontWeight,
  },
});
