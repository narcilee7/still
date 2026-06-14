import React, { useCallback, useState } from 'react';
import { ScrollView, StyleSheet, Text, TextInput, View } from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { SafeAreaView } from 'react-native-safe-area-context';
import { RootStackParamList } from '../navigation/types';
import { useStore } from '../store/useStore';
import { updatePost } from '../services/postApi';
import { useTranslation } from 'react-i18next';
import { colors, spacing, typography, QuietButton } from '@still/design-system';

type Props = NativeStackScreenProps<RootStackParamList, 'EditPost'>;

export function EditPostScreen({ route, navigation }: Props) {
  const { t } = useTranslation();
  const { postId } = route.params;
  const post = useStore((state) => state.posts.find((p) => p.id === postId));
  const updatePostInStore = useStore((state) => state.updatePost);
  const [title, setTitle] = useState(post?.title ?? '');
  const [description, setDescription] = useState(post?.description ?? '');
  const [saving, setSaving] = useState(false);

  const save = useCallback(async () => {
    if (!post) return;
    setSaving(true);
    try {
      const updated = await updatePost({
        id: postId,
        mood: post.mood,
        title: title.trim(),
        description: description.trim(),
      });
      updatePostInStore(postId, {
        title: updated.title,
        description: updated.description,
        mood: updated.mood,
      });
      navigation.goBack();
    } catch (err) {
      console.error('update post failed', err);
    } finally {
      setSaving(false);
    }
  }, [post, postId, title, description, updatePostInStore, navigation]);

  if (!post) {
    return (
      <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
        <View style={styles.centered}>
          <Text style={styles.emptyText}>{t('editPost.notFound')}</Text>
          <QuietButton title={t('editPost.goBack')} onPress={() => navigation.goBack()} />
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <ScrollView contentContainerStyle={styles.content}>
        <Text style={styles.label}>{t('editPost.title')}</Text>
        <TextInput
          value={title}
          onChangeText={setTitle}
          style={styles.input}
          placeholder={t('create.edit.titlePlaceholder')}
          placeholderTextColor={colors.secondary}
          maxLength={80}
        />
        <Text style={styles.label}>{t('editPost.description')}</Text>
        <TextInput
          value={description}
          onChangeText={setDescription}
          style={[styles.input, styles.inputMultiline]}
          placeholder={t('create.edit.descriptionPlaceholder')}
          placeholderTextColor={colors.secondary}
          multiline
          maxLength={240}
        />
      </ScrollView>
      <View style={styles.footer}>
        <QuietButton title={t('editPost.save')} onPress={save} disabled={saving} />
      </View>
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
  },
  emptyText: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.secondary,
    marginBottom: spacing.md,
  },
  content: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  label: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    marginTop: spacing.md,
    marginBottom: spacing.sm,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  input: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.primary,
    backgroundColor: colors.white,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  inputMultiline: {
    minHeight: 96,
    textAlignVertical: 'top',
  },
  footer: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.md,
    paddingBottom: spacing.lg,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
});
