import React, { useState } from "react";
import {
  Dimensions,
  Image,
  Pressable,
  StyleSheet,
  Text,
  View,
  ViewStyle,
} from "react-native";
import { Post } from "@still/shared-types";
import { colors, spacing, typography } from "../theme";
import { MoodTag } from "./MoodTag";
import { ResonateButton } from "./ResonateButton";

const { height: screenHeight } = Dimensions.get("window");

export interface PostCardProps {
  post: Post;
  variant: "full" | "compact";
  resonated?: boolean;
  onResonate?: () => void;
  onShare?: () => void;
  style?: ViewStyle;
}

export function PostCard({
  post,
  variant,
  resonated,
  onResonate,
  onShare,
  style,
}: PostCardProps) {
  const [imageLoaded, setImageLoaded] = useState(false);

  if (variant === "compact") {
    return (
      <View style={[compactStyles.root, style]}>
        <Image
          source={{ uri: post.imageUrl }}
          style={compactStyles.image}
          resizeMode="cover"
          onLoad={() => setImageLoaded(true)}
        />
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
      </View>
    );
  }

  return (
    <View style={[fullStyles.root, style]}>
      <View style={fullStyles.imageContainer}>
        <Image
          source={{ uri: post.imageUrl }}
          style={[fullStyles.image, imageLoaded && fullStyles.imageVisible]}
          resizeMode="cover"
          onLoad={() => setImageLoaded(true)}
        />
        {!imageLoaded && <View style={fullStyles.imagePlaceholder} />}
      </View>

      <View style={fullStyles.textArea}>
        <MoodTag mood={post.mood} />
        <Text style={fullStyles.title}>{post.title}</Text>
        <Text style={fullStyles.description}>{post.description}</Text>

        <View style={fullStyles.actions}>
          <ResonateButton
            count={post.resonanceCount}
            resonated={resonated}
            onPress={onResonate}
          />
          {onShare && (
            <Pressable
              onPress={onShare}
              style={fullStyles.share}
              accessibilityRole="button"
            >
              <Text style={fullStyles.shareLabel}>Share</Text>
            </Pressable>
          )}
        </View>
      </View>
    </View>
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
    opacity: 0,
  },
  imageVisible: {
    opacity: 1,
  },
  imagePlaceholder: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: colors.border,
  },
  textArea: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
    paddingBottom: spacing.xxl,
    minHeight: screenHeight * 0.32,
    justifyContent: "flex-end",
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
    flexDirection: "row",
    alignItems: "center",
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
    flexDirection: "row",
    backgroundColor: colors.background,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    gap: spacing.md,
  },
  image: {
    width: 88,
    height: 110,
    borderRadius: 8,
    backgroundColor: colors.border,
  },
  body: {
    flex: 1,
    justifyContent: "center",
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
