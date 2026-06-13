import React, { useCallback, useEffect, useState } from 'react';
import { Dimensions, FlatList, ListRenderItem, RefreshControl, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Post } from '@still/shared-types';
import { PostCard, colors } from '@still/design-system';
import { listFeed, resonate } from '../services/postApi';
import { useStore } from '../store/useStore';

const { height: WINDOW_HEIGHT } = Dimensions.get('window');

export function FeedScreen() {
  const posts = useStore((state) => state.posts);
  const setPosts = useStore((state) => state.setPosts);
  const resonatedIds = useStore((state) => state.resonatedPostIds);
  const storeToggleResonate = useStore((state) => state.toggleResonate);
  const updatePost = useStore((state) => state.updatePost);
  const [listHeight, setListHeight] = useState(WINDOW_HEIGHT);
  const [refreshing, setRefreshing] = useState(false);

  const load = useCallback(async () => {
    try {
      const data = await listFeed();
      setPosts(data);
    } catch (err) {
      console.error('feed load failed', err);
    } finally {
      setRefreshing(false);
    }
  }, [setPosts]);

  useEffect(() => {
    load();
  }, [load]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    load();
  }, [load]);

  const handleResonate = useCallback(
    async (postId: string) => {
      storeToggleResonate(postId);
      try {
        const updated = await resonate(postId);
        updatePost(postId, { resonanceCount: updated.resonanceCount });
      } catch (err) {
        console.error('resonate failed', err);
      }
    },
    [storeToggleResonate, updatePost]
  );

  const renderItem: ListRenderItem<Post> = useCallback(
    ({ item }) => (
      <PostCard
        post={item}
        variant="full"
        resonated={resonatedIds.has(item.id)}
        onResonate={() => handleResonate(item.id)}
        onShare={() => {}}
        style={{ height: listHeight }}
      />
    ),
    [listHeight, resonatedIds, handleResonate]
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
          refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor={colors.secondary} />}
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
});
