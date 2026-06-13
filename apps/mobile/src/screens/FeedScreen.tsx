import React, { useCallback, useEffect, useState } from 'react';
import {
  Dimensions,
  FlatList,
  ListRenderItem,
  RefreshControl,
  Share,
  StyleSheet,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Post } from '@still/shared-types';
import { EmptyState, ErrorState, LoadingSpinner, PostCard, colors } from '@still/design-system';
import { listFeed, resonate } from '../services/postApi';
import { useStore } from '../store/useStore';

const { height: WINDOW_HEIGHT } = Dimensions.get('window');

type LoadState = 'idle' | 'loading' | 'error';

export function FeedScreen() {
  const posts = useStore((state) => state.posts);
  const setPosts = useStore((state) => state.setPosts);
  const appendPosts = useStore((state) => state.appendPosts);
  const resonatedIds = useStore((state) => state.resonatedPostIds);
  const setResonated = useStore((state) => state.setResonated);
  const updatePost = useStore((state) => state.updatePost);
  const [listHeight, setListHeight] = useState(WINDOW_HEIGHT);
  const [refreshing, setRefreshing] = useState(false);
  const [loadState, setLoadState] = useState<LoadState>('loading');
  const [nextPageToken, setNextPageToken] = useState('');
  const [loadingMore, setLoadingMore] = useState(false);

  const load = useCallback(
    async (mode: 'refresh' | 'append' = 'refresh') => {
      if (mode === 'refresh') {
        setLoadState('loading');
      } else {
        setLoadingMore(true);
      }
      try {
        const token = mode === 'append' ? nextPageToken : '';
        const page = await listFeed(token);
        if (mode === 'refresh') {
          setPosts(page.posts);
        } else {
          appendPosts(page.posts);
        }
        setNextPageToken(page.nextPageToken);
        setLoadState('idle');
      } catch (err) {
        console.error('feed load failed', err);
        if (mode === 'refresh') {
          setLoadState('error');
        }
      } finally {
        setRefreshing(false);
        setLoadingMore(false);
      }
    },
    [setPosts, appendPosts, nextPageToken]
  );

  useEffect(() => {
    load('refresh');
  }, [load]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    load('refresh');
  }, [load]);

  const onEndReached = useCallback(() => {
    if (loadingMore || !nextPageToken) return;
    load('append');
  }, [load, loadingMore, nextPageToken]);

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

  const handleShare = useCallback(async (post: Post) => {
    try {
      await Share.share({
        message: `${post.title}\n\n${post.description}\n\n${post.imageUrl}`,
      });
    } catch (err) {
      console.error('share failed', err);
    }
  }, []);

  const renderItem: ListRenderItem<Post> = useCallback(
    ({ item }) => (
      <PostCard
        post={item}
        variant="full"
        resonated={resonatedIds.has(item.id)}
        onResonate={() => handleResonate(item.id)}
        onShare={() => handleShare(item)}
        style={{ height: listHeight }}
      />
    ),
    [listHeight, resonatedIds, handleResonate, handleShare]
  );

  const keyExtractor = useCallback((item: Post) => item.id, []);

  const getItemLayout = useCallback(
    (_data: ArrayLike<Post> | null | undefined, index: number) => ({
      length: listHeight,
      offset: listHeight * index,
      index,
    }),
    [listHeight]
  );

  const renderFooter = useCallback(() => {
    if (!loadingMore) return null;
    return <LoadingSpinner size="small" />;
  }, [loadingMore]);

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
            title="Could not load feed"
            message="Something went wrong while fetching moments. Please try again."
            onRetry={() => load('refresh')}
          />
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <View style={styles.list} onLayout={(e) => setListHeight(e.nativeEvent.layout.height)}>
        <FlatList
          data={posts}
          renderItem={renderItem}
          keyExtractor={keyExtractor}
          pagingEnabled
          showsVerticalScrollIndicator={false}
          getItemLayout={getItemLayout}
          decelerationRate="fast"
          onEndReached={onEndReached}
          onEndReachedThreshold={0.5}
          ListFooterComponent={renderFooter}
          ListEmptyComponent={
            <EmptyState title="No moments yet" subtitle="Be the first to share a quiet moment." />
          }
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={onRefresh}
              tintColor={colors.secondary}
            />
          }
        />
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  list: {
    flex: 1,
  },
  centered: {
    flex: 1,
    justifyContent: 'center',
  },
});
