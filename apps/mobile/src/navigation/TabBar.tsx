import { Pressable, StyleSheet, Text, View } from 'react-native';
import { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useTranslation } from 'react-i18next';
import { colors, spacing, typography } from '@still/design-system';

export function TabBar({ state, descriptors, navigation }: BottomTabBarProps) {
  const { t } = useTranslation();

  const labelFor = (name: string) => {
    switch (name) {
      case 'Feed':
        return t('tabBar.feed');
      case 'Create':
        return t('tabBar.create');
      case 'Profile':
        return t('tabBar.profile');
      default:
        return name;
    }
  };
  return (
    <SafeAreaView edges={['bottom']} style={styles.container}>
      <View style={styles.bar}>
        {state.routes.map((route, index) => {
          const { options } = descriptors[route.key];
          const rawLabel =
            options.tabBarLabel ?? options.title ?? labelFor(route.name) ?? route.name;
          const label = typeof rawLabel === 'string' ? rawLabel : route.name;
          const isFocused = state.index === index;

          const onPress = () => {
            const event = navigation.emit({
              type: 'tabPress',
              target: route.key,
              canPreventDefault: true,
            });

            if (!isFocused && !event.defaultPrevented) {
              navigation.navigate(route.name);
            }
          };

          return (
            <Pressable
              key={route.key}
              onPress={onPress}
              style={styles.tab}
              accessibilityRole="button"
              accessibilityState={{ selected: isFocused }}
              accessibilityLabel={label}
            >
              <Text style={[styles.label, isFocused && styles.labelFocused]}>{label}</Text>
            </Pressable>
          );
        })}
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.background,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  bar: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    paddingTop: spacing.md,
    paddingBottom: spacing.md,
  },
  tab: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  label: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.secondary,
  },
  labelFocused: {
    color: colors.primary,
  },
});
