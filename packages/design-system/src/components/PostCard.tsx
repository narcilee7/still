import React, { useState } from 'react';
import { Dimensions, Pressable, StyleSheet, Text, View, ViewStyle } from 'react-native';
import { Image } from 'expo-image';
import { Post } from '@still/shared-types';
import { colors, spacing, typography } from '../theme';
import { MoodTag } from './MoodTag'
import { ResonateButton } from './ResonateButton';

const { height: screenHeight } = Dimensions.get('window');

export interface PostCardProps {
  post: Post;
  variant: 'full' | 'compact';
  resonated?: boolean;
  onResonate?: () => void;
  onShare?: () => void;
  onPress?: () => void;
  style?: ViewStyle;
}

export function PostCard({
  post,
  variant,
  resonated,
  onResonate,
  onShare,
  onPress,
  style,
}: PostCardProps) {
  const [imageError, setImageError] = useState(false);

  if (variant === 'compact') {
    return (
      <Pressable
        onPress={onPress}
        disabled={!onPress}
        style={[compactStyles.root, style]}
      >
        <View style={compactStyles.imageContainer}>
          <Image
            source={{ uri: post.imageUrl }}
            style={compactStyles.image}
            contentFit="cover"
            cachePolicy="memory-disk"
            transition={500}
            onError={() => setImageError(true)}
          />
          {imageError && (
            <View style={compactStyles.imagePlaceholder}>
              <Text style={compactStyles.imagePlaceholderText}>?</Text>
            </View>
          )}
        </View>
        <View style={compactStyles.body}>
          <MoodTag mood={post.mood} variant="small" />
          <Text style={compactStyles.title} numberOfLines={1}>
            {post.title}
          </Text>
          <Text style={compactStyles.description} numberOfLines={2}>
            {post.description}
          </Text>
          <View style={compactStyles.footer}>
            <ResonateButton
              count={post.resonanceCount}
              resonated={resonated}
              onPress={onResonate}
            />
          </View>
        </View>
      </Pressable>
    );
  }

  return (
    <Pressable
      onPress={onPress}
      disabled={!onPress}
      style={[fullStyles.root, style]}
    >
      <View style={fullStyles.imageContainer}>
        <Image
          source={{ uri: post.imageUrl }}
          style={fullStyles.image}
          contentFit="cover"
          cachePolicy="memory-disk"
          transition={500}
          onError={() => setImageError(true)}
        />
        {imageError && (
          <View style={fullStyles.imagePlaceholder}>
            <Text style={fullStyles.imagePlaceholderText}>?</Text>
          </View>
        )}
      </View>

      <View style={fullStyles.textArea}>
        <MoodTag mood={post.mood} />
        <Text style={fullStyles.title}>{post.title}</Text>
        <Text style={fullStyles.description}>{post.description}</Text>

        <View style={fullStyles.actions}>
          <ResonateButton count={post.resonanceCount} resonated={resonated} onPress={onResonate} />
          {onShare && (
            <Pressable onPress={onShare} style={fullStyles.share} accessibilityRole="button">
              <Text style={fullStyles.shareLabel}>Share</Text>
            </Pressable>
          )}
        </View>
      </View>
    </Pressable>
  );
}

const fullStyles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: colors.background,
  },
  imageContainer: {
    flex: 1,
    backgroundColor: colors.border,
  },
  image: {
    flex: 1,
  },
  imagePlaceholder: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  imagePlaceholderText: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.secondary,
  },
  textArea: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
    paddingBottom: spacing.xxl,
    minHeight: screenHeight * 0.32,
    justifyContent: 'flex-end',
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
  actions: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: spacing.lg,
    gap: spacing.md,
  },
  share: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
  },
  shareLabel: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
  },
});

const compactStyles = StyleSheet.create({
  root: {
    flexDirection: 'row',
    backgroundColor: colors.background,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    gap: spacing.md,
  },
  imageContainer: {
    width: 88,
    height: 110,
    borderRadius: 8,
    backgroundColor: colors.border,
    overflow: 'hidden',
  },
  image: {
    width: 88,
    height: 110,
    backgroundColor: colors.border,
  },
  imagePlaceholder: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  imagePlaceholderText: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.secondary,
  },
  body: {
    flex: 1,
    justifyContent: 'center',
  },
  title: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.primary,
    marginTop: spacing.xs,
  },
  description: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    marginTop: spacing.xs,
  },
  footer: {
    marginTop: spacing.sm,
  },
});
