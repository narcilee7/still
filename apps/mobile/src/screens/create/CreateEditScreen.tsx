import React, { useCallback, useEffect, useState } from 'react';
import {
  Image,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Mood, MOODS } from '@still/shared-types';
import { colors, spacing, typography, ErrorState, QuietButton } from '@still/design-system';
import { CreateStackParamList } from '../../navigation/types';
import { analyzeImage, createPost, getUploadURL, uploadImage } from '../../services/postApi';
import { useStore } from '../../store/useStore';

type Props = NativeStackScreenProps<CreateStackParamList, 'CreateEdit'>;

const LOADING_TEXTS = ['Finding words…', 'Looking closer…', 'Finding words…'];

export function CreateEditScreen({ route, navigation }: Props) {
  const { imageUri } = route.params;
  const addPost = useStore((state) => state.addPost);

  const [step, setStep] = useState<'uploading' | 'analyzing' | 'editing' | 'publishing' | 'error'>('uploading');
  const [loadingText, setLoadingText] = useState(LOADING_TEXTS[0]);
  const [publicUrl, setPublicUrl] = useState('');
  const [mood, setMood] = useState<Mood>('still');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [publishError, setPublishError] = useState<string | null>(null);

  const prepare = useCallback(async () => {
    let textInterval: ReturnType<typeof setInterval> | null = null;
    setStep('uploading');
    setPublishError(null);

    try {
      textInterval = setInterval(() => {
        setLoadingText((prev) => {
          const idx = LOADING_TEXTS.indexOf(prev);
          return LOADING_TEXTS[(idx + 1) % LOADING_TEXTS.length];
        });
      }, 900);

      const blob = await fetch(imageUri).then((r) => r.blob());
      const contentType = blob.type || 'image/jpeg';
      const filename = imageUri.split('/').pop() || 'image.jpg';

      const { uploadUrl, publicUrl: pub } = await getUploadURL(filename, contentType);
      await uploadImage(imageUri, uploadUrl, contentType);
      setPublicUrl(pub);

      setStep('analyzing');
      const result = await analyzeImage(pub);
      setMood(result.mood);
      setTitle(result.title);
      setDescription(result.description);
      setStep('editing');
    } catch (err) {
      console.error('prepare failed', err);
      setStep('error');
    } finally {
      if (textInterval) clearInterval(textInterval);
    }
  }, [imageUri]);

  useEffect(() => {
    prepare();
  }, [prepare]);

  const publish = useCallback(async () => {
    if (!publicUrl) return;
    setStep('publishing');
    setPublishError(null);
    try {
      const post = await createPost({
        imageUrl: publicUrl,
        mood,
        title: title.trim(),
        description: description.trim(),
      });
      addPost(post);
      navigation.replace('CreateSuccess', { postId: post.id });
    } catch (err) {
      console.error('publish failed', err);
      setPublishError('We could not publish your moment. Please try again.');
      setStep('editing');
    }
  }, [publicUrl, mood, title, description, addPost, navigation]);

  if (step === 'uploading' || step === 'analyzing' || step === 'publishing') {
    return (
      <SafeAreaView edges={['top', 'left', 'right']} style={styles.centered}>
        <Text style={styles.loadingText}>
          {step === 'uploading' ? 'Uploading…' : step === 'publishing' ? 'Publishing…' : loadingText}
        </Text>
      </SafeAreaView>
    );
  }

  if (step === 'error') {
    return (
      <SafeAreaView edges={['top', 'left', 'right']} style={styles.centered}>
        <ErrorState
          title="Could not prepare your moment"
          message="Something went wrong while uploading or analyzing your photo. Please try again."
          onRetry={prepare}
        />
      </SafeAreaView>
    );
  }

  const canPublish = title.trim().length > 0 && description.trim().length > 0;

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <ScrollView
        style={styles.scroll}
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
      >
        {publicUrl ? (
          <Image source={{ uri: publicUrl }} style={styles.preview} resizeMode="cover" />
        ) : (
          <Image source={{ uri: imageUri }} style={styles.preview} resizeMode="cover" />
        )}

        <Text style={styles.sectionLabel}>Mood</Text>
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={styles.moodList}
        >
          {MOODS.map((m) => {
            const selected = m === mood;
            return (
              <Pressable
                key={m}
                onPress={() => setMood(m)}
                style={[styles.moodChip, selected && styles.moodChipSelected]}
                accessibilityRole="radio"
                accessibilityState={{ selected }}
              >
                <Text style={[styles.moodChipText, selected && styles.moodChipTextSelected]}>
                  {m}
                </Text>
              </Pressable>
            );
          })}
        </ScrollView>

        <Text style={styles.sectionLabel}>Title</Text>
        <TextInput
          value={title}
          onChangeText={setTitle}
          style={styles.input}
          placeholder="A few words"
          placeholderTextColor={colors.secondary}
          maxLength={80}
        />

        <Text style={styles.sectionLabel}>Description</Text>
        <TextInput
          value={description}
          onChangeText={setDescription}
          style={[styles.input, styles.inputMultiline]}
          placeholder="What lingers?"
          placeholderTextColor={colors.secondary}
          multiline
          maxLength={240}
        />
      </ScrollView>

      <View style={styles.footer}>
        {publishError ? <Text style={styles.publishError}>{publishError}</Text> : null}
        <QuietButton title="Publish" onPress={publish} disabled={!canPublish} />
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
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: spacing.lg,
  },
  scroll: {
    flex: 1,
  },
  scrollContent: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  loadingText: {
    fontSize: typography.title.fontSize,
    lineHeight: typography.title.lineHeight,
    color: colors.secondary,
  },
  preview: {
    width: '100%',
    height: 280,
    borderRadius: 12,
    backgroundColor: colors.border,
    marginVertical: spacing.lg,
  },
  sectionLabel: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    marginTop: spacing.md,
    marginBottom: spacing.sm,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  moodList: {
    gap: spacing.sm,
    paddingBottom: spacing.sm,
  },
  moodChip: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.white,
  },
  moodChipSelected: {
    backgroundColor: colors.primary,
    borderColor: colors.primary,
  },
  moodChipText: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.primary,
  },
  moodChipTextSelected: {
    color: colors.white,
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
  publishError: {
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.primary,
    textAlign: 'center',
    marginBottom: spacing.md,
  },
});
