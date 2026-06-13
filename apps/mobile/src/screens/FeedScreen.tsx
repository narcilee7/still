import React, { useCallback, useState } from 'react';
import { Dimensions, FlatList, ListRenderItem, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Post } from '@still/shared-types';
import { PostCard, colors } from '@still/design-system';
import { useStore } from '../store/useStore';

const { height: WINDOW_HEIGHT } = Dimensions.get('window');

export function FeedScreen() {
  const posts = useStore((state) => state.posts);
  const resonatedIds = useStore((state) => state.resonatedPostIds);
  const toggleResonate = useStore((state) => state.toggleResonate);
  const [listHeight, setListHeight] = useState(WINDOW_HEIGHT);

  const renderItem: ListRenderItem<Post> = useCallback(
    ({ item }) => (
      <PostCard
        post={item}
        variant="full"
        resonated={resonatedIds.has(item.id)}
        onResonate={() => toggleResonate(item.id)}
        onShare={() => {}}
        style={{ height: listHeight }}
      />
    ),
    [listHeight, resonatedIds, toggleResonate]
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
