import React, { useCallback } from 'react';
import { Alert, Pressable, ScrollView, Share, StyleSheet, Text, View } from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Image } from 'expo-image';
import { RootStackParamList } from '../navigation/types';
import { useStore } from '../store/useStore';
import { deletePost } from '../services/postApi';
import { colors, spacing, typography, QuietButton } from '@still/design-system';

type Props = NativeStackScreenProps<RootStackParamList, 'PostDetail'>;

export function PostDetailScreen({ route, navigation }: Props) {
  const { postId } = route.params;
  const post = useStore((state) => state.posts.find((p) => p.id === postId));
  const user = useStore((state) => state.user);
  const removePost = useStore((state) => state.removePost);
  const isOwner = post?.userId === user.id;

  const handleShare = useCallback(async () => {
    if (!post) return;
    try {
      await Share.share({
        message: `${post.title}\n\n${post.description}\n\n${post.imageUrl}`,
      });
    } catch (err) {
      console.error('share failed', err);
    }
  }, [post]);

  const handleDelete = useCallback(() => {
    if (!post) return;
    Alert.alert('Delete moment?', 'This cannot be undone.', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: async () => {
          try {
            await deletePost(post.id);
            removePost(post.id);
            navigation.goBack();
          } catch (err) {
            console.error('delete post failed', err);
          }
        },
      },
    ]);
  }, [post, removePost, navigation]);

  if (!post) {
    return (
      <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
        <View style={styles.centered}>
          <Text style={styles.emptyText}>Moment not found</Text>
          <QuietButton title="Go back" onPress={() => navigation.goBack()} />
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <ScrollView contentContainerStyle={styles.content}>
        <Pressable onPress={() => navigation.goBack()} style={styles.close}>
          <Text style={styles.closeText}>Close</Text>
        </Pressable>
        <Image
          source={{ uri: post.imageUrl }}
          style={styles.image}
          contentFit="cover"
          cachePolicy="memory-disk"
          transition={500}
        />
        <View style={styles.meta}>
          <Text style={styles.mood}>{post.mood}</Text>
          <Text style={styles.title}>{post.title}</Text>
          <Text style={styles.description}>{post.description}</Text>
          <Text style={styles.count}>{post.resonanceCount} resonances</Text>
        </View>
        <View style={styles.actions}>
          <QuietButton title="Share" variant="secondary" onPress={handleShare} />
          {isOwner && (
            <>
              <QuietButton
                title="Edit"
                variant="secondary"
                onPress={() => navigation.navigate('EditPost', { postId: post.id })}
              />
              <QuietButton title="Delete" variant="secondary" onPress={handleDelete} />
            </>
          )}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  centered: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
  },
  emptyText: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.secondary,
    marginBottom: spacing.md,
  },
  content: {
    paddingBottom: spacing.xl,
  },
  close: {
    alignSelf: 'flex-start',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  closeText: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.secondary,
  },
  image: {
    width: '100%',
    height: 420,
    backgroundColor: colors.border,
  },
  meta: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
  },
  mood: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  title: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.primary,
    marginTop: spacing.sm,
  },
  description: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.secondary,
    marginTop: spacing.sm,
  },
  count: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    marginTop: spacing.lg,
  },
  actions: {
    flexDirection: 'row',
    gap: spacing.md,
    paddingHorizontal: spacing.lg,
    marginTop: spacing.xl,
  },
});
