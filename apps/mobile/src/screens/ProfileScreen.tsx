import React, { useCallback } from 'react';
import { FlatList, ListRenderItem, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Post } from '@still/shared-types';
import { PostCard, colors, spacing, typography } from '@still/design-system';
import { useStore } from '../store/useStore';

export function ProfileScreen() {
  const user = useStore((state) => state.user);
  const posts = useStore((state) => state.posts.filter((p) => p.userId === state.user.id));
  const resonatedIds = useStore((state) => state.resonatedPostIds);
  const toggleResonate = useStore((state) => state.toggleResonate);

  const renderItem: ListRenderItem<Post> = useCallback(
    ({ item }) => (
      <PostCard
        post={item}
        variant="compact"
        resonated={resonatedIds.has(item.id)}
        onResonate={() => toggleResonate(item.id)}
      />
    ),
    [resonatedIds, toggleResonate]
  );

  const keyExtractor = useCallback((item: Post) => item.id, []);

  const Header = () => (
    <View style={styles.header}>
      <View style={styles.avatar}>
        <Text style={styles.avatarText}>{user.username.charAt(0).toUpperCase()}</Text>
      </View>
      <Text style={styles.username}>@{user.username}</Text>
      <View style={styles.stats}>
        <View style={styles.stat}>
          <Text style={styles.statNumber}>{user.postsCount}</Text>
          <Text style={styles.statLabel}>Moments</Text>
        </View>
        <View style={styles.stat}>
          <Text style={styles.statNumber}>{user.resonancesCount}</Text>
          <Text style={styles.statLabel}>Resonances</Text>
        </View>
      </View>
    </View>
  );

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <FlatList
        data={posts}
        renderItem={renderItem}
        keyExtractor={keyExtractor}
        ListHeaderComponent={Header}
        contentContainerStyle={styles.listContent}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  listContent: {
    paddingBottom: 32,
  },
  header: {
    alignItems: 'center',
    paddingVertical: 40,
    paddingHorizontal: spacing.lg,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: spacing.md,
  },
  avatarText: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.primary,
  },
  username: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.primary,
    marginBottom: spacing.lg,
  },
  stats: {
    flexDirection: 'row',
    gap: spacing.xl,
  },
  stat: {
    alignItems: 'center',
  },
  statNumber: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.primary,
  },
  statLabel: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    marginTop: spacing.xs,
  },
});
