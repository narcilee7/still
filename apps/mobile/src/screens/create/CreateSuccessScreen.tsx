import React, { useCallback } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, typography, QuietButton } from '@still/design-system';
import { CreateStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<CreateStackParamList, 'CreateSuccess'>;

export function CreateSuccessScreen({ navigation }: Props) {
  const addAnother = useCallback(() => {
    navigation.popToTop();
  }, [navigation]);

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <View style={styles.content}>
        <Text style={styles.title}>Your moment lingers.</Text>
        <Text style={styles.subtitle}>It is now part of the quiet stream.</Text>
      </View>
      <View style={styles.footer}>
        <QuietButton title="Add another" onPress={addAnother} />
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    paddingHorizontal: spacing.lg,
  },
  content: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: typography.mood.fontSize,
    lineHeight: typography.mood.lineHeight,
    color: colors.primary,
    marginBottom: spacing.sm,
    textAlign: 'center',
  },
  subtitle: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.secondary,
    textAlign: 'center',
  },
  footer: {
    paddingBottom: spacing.lg,
  },
});
