import React, { useCallback, useEffect, useState } from 'react';
import { FlatList, ListRenderItem, RefreshControl, StyleSheet, Text, View } from 'react-native';
import { Image } from 'expo-image';
import { useAuth } from '@clerk/expo';
import { useNavigation } from '@react-navigation/native';
import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Post } from '@still/shared-types';
import {
  EmptyState,
  ErrorState,
  LoadingSpinner,
  PostCard,
  QuietButton,
  colors,
  spacing,
  typography,
} from '@still/design-system';
import { getProfile, listFeed, resonate } from '../services/postApi';
import { RootStackParamList } from '../navigation/types';
import { useTranslation } from 'react-i18next';
import { useShallow } from 'zustand/react/shallow';
import { useStore } from '../store/useStore';

type LoadState = 'idle' | 'loading' | 'error';

export function ProfileScreen() {
  const { t, i18n } = useTranslation();
  const navigation = useNavigation<NativeStackNavigationProp<RootStackParamList>>();
  const user = useStore((state) => state.user);
  const posts = useStore(
    useShallow((state) => state.posts.filter((p) => p.userId === state.user.id))
  );
  const setUser = useStore((state) => state.setUser);
  const setPosts = useStore((state) => state.setPosts);
  const resonatedIds = useStore((state) => state.resonatedPostIds);
  const setResonated = useStore((state) => state.setResonated);
  const updatePost = useStore((state) => state.updatePost);
  const [loadState, setLoadState] = useState<LoadState>('loading');
  const [refreshing, setRefreshing] = useState(false);
  const { signOut } = useAuth();

  const load = useCallback(async () => {
    try {
      const [profile, feed] = await Promise.all([getProfile(user.id), listFeed()]);
      setUser((prev) => ({
        ...prev,
        username: profile.user.username,
        avatarUrl: profile.user.avatarUrl,
        postsCount: profile.postsCount,
        resonancesCount: profile.resonancesCount,
      }));
      setPosts(feed.posts);
      setLoadState('idle');
    } catch (err) {
      console.error('profile load failed', err);
      setLoadState('error');
    }
  }, [user.id, setUser, setPosts]);

  useEffect(() => {
    load();
  }, [load]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    load().finally(() => setRefreshing(false));
  }, [load]);

  const handleResonate = useCallback(
    async (postId: string) => {
      const wasResonated = resonatedIds.has(postId);
      setResonated(postId, !wasResonated);
      try {
        const result = await resonate(postId);
        updatePost(postId, { resonanceCount: result.post.resonanceCount });
        setResonated(postId, result.hasResonated);
      } catch (err) {
        console.error('resonate failed', err);
        setResonated(postId, wasResonated);
      }
    },
    [resonatedIds, setResonated, updatePost]
  );

  const renderItem: ListRenderItem<Post> = useCallback(
    ({ item }) => (
      <PostCard
        post={item}
        variant="compact"
        resonated={resonatedIds.has(item.id)}
        onResonate={() => handleResonate(item.id)}
        onPress={() => navigation.navigate('PostDetail', { postId: item.id })}
      />
    ),
    [resonatedIds, handleResonate, navigation]
  );

  const keyExtractor = useCallback((item: Post) => item.id, []);

  const handleSignOut = useCallback(async () => {
    try {
      await signOut();
    } catch (err) {
      console.error('sign out failed', err);
    }
  }, [signOut]);

  const toggleLanguage = useCallback(() => {
    i18n.changeLanguage(i18n.language.startsWith('zh') ? 'en' : 'zh');
  }, [i18n]);

  const Header = () => (
    <View style={styles.header}>
      <View style={styles.avatar}>
        {user.avatarUrl ? (
          <Image source={{ uri: user.avatarUrl }} style={styles.avatarImage} />
        ) : (
          <Text style={styles.avatarText}>{user.username.charAt(0).toUpperCase()}</Text>
        )}
      </View>
      <Text style={styles.username}>@{user.username}</Text>
      <View style={styles.stats}>
        <View style={styles.stat}>
          <Text style={styles.statNumber}>{user.postsCount}</Text>
          <Text style={styles.statLabel}>{t('profile.moments')}</Text>
        </View>
        <View style={styles.stat}>
          <Text style={styles.statNumber}>{user.resonancesCount}</Text>
          <Text style={styles.statLabel}>{t('profile.resonances')}</Text>
        </View>
      </View>
      <View style={styles.signOut}>
        <QuietButton title={t('profile.signOut')} variant="secondary" onPress={handleSignOut} />
      </View>
      <View style={styles.languageToggle}>
        <QuietButton
          title={t('profile.switchLanguage')}
          variant="secondary"
          onPress={toggleLanguage}
        />
      </View>
    </View>
  );

  if (loadState === 'loading' && posts.length === 0) {
    return (
      <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
        <View style={styles.centered}>
          <LoadingSpinner size="large" />
        </View>
      </SafeAreaView>
    );
  }

  if (loadState === 'error' && posts.length === 0) {
    return (
      <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
        <View style={styles.centered}>
          <ErrorState
            title={t('profile.errorTitle')}
            message={t('profile.errorMessage')}
            onRetry={load}
          />
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <FlatList
        data={posts}
        renderItem={renderItem}
        keyExtractor={keyExtractor}
        ListHeaderComponent={Header}
        contentContainerStyle={styles.listContent}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={colors.secondary}
          />
        }
        ListEmptyComponent={
          <EmptyState title={t('profile.emptyTitle')} subtitle={t('profile.emptySubtitle')} />
        }
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
  centered: {
    flex: 1,
    justifyContent: 'center',
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
    overflow: 'hidden',
  },
  avatarImage: {
    width: 80,
    height: 80,
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
  signOut: {
    marginTop: spacing.xl,
    minWidth: 160,
  },
  languageToggle: {
    marginTop: spacing.md,
    minWidth: 160,
  },
});
