import React, { useCallback, useEffect, useState } from 'react';
import { Linking, StyleSheet, Text, View } from 'react-native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import * as ImagePicker from 'expo-image-picker';
import { SafeAreaView } from 'react-native-safe-area-context';
import { colors, spacing, typography, QuietButton } from '@still/design-system';
import { CreateStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<CreateStackParamList, 'CreateSelect'>;

export function CreateSelectScreen({ navigation }: Props) {
  const [mediaPermissionStatus, setMediaPermissionStatus] = useState<ImagePicker.PermissionStatus | null>(null);
  const [cameraPermissionStatus, setCameraPermissionStatus] = useState<ImagePicker.PermissionStatus | null>(null);

  useEffect(() => {
    ImagePicker.requestMediaLibraryPermissionsAsync().then((result) => {
      setMediaPermissionStatus(result.status);
    });
    ImagePicker.requestCameraPermissionsAsync().then((result) => {
      setCameraPermissionStatus(result.status);
    });
  }, []);

  const openLibrary = useCallback(async () => {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ImagePicker.MediaTypeOptions.Images,
      allowsEditing: true,
      aspect: [3, 4],
      quality: 0.85,
    });

    if (!result.canceled && result.assets.length > 0) {
      navigation.navigate('CreateEdit', { imageUri: result.assets[0].uri });
    }
  }, [navigation]);

  const takePhoto = useCallback(async () => {
    const result = await ImagePicker.launchCameraAsync({
      allowsEditing: true,
      aspect: [3, 4],
      quality: 0.85,
    });

    if (!result.canceled && result.assets.length > 0) {
      navigation.navigate('CreateEdit', { imageUri: result.assets[0].uri });
    }
  }, [navigation]);

  const openSettings = useCallback(() => {
    Linking.openSettings();
  }, []);

  const showMediaPermissionHint = mediaPermissionStatus === ImagePicker.PermissionStatus.DENIED;
  const showCameraPermissionHint = cameraPermissionStatus === ImagePicker.PermissionStatus.DENIED;

  return (
    <SafeAreaView edges={['top', 'left', 'right']} style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Add a moment</Text>
        <Text style={styles.subtitle}>A photo, a feeling, a few quiet words.</Text>
      </View>

      <View style={styles.actions}>
        <QuietButton title="Choose from library" onPress={openLibrary} variant="primary" />
        <QuietButton
          title="Take photo"
          onPress={takePhoto}
          variant="secondary"
        />
      </View>

      {showMediaPermissionHint && (
        <Text style={styles.permissionHint}>
          Gallery access is needed to choose a photo.
        </Text>
      )}
      {showCameraPermissionHint && (
        <View style={styles.permissionRow}>
          <Text style={styles.permissionHint}>
            Camera access is needed to take a photo.
          </Text>
          <QuietButton title="Open Settings" onPress={openSettings} variant="secondary" />
        </View>
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    paddingHorizontal: spacing.lg,
    justifyContent: 'center',
  },
  header: {
    marginBottom: spacing.xxl,
    alignItems: 'center',
  },
  title: {
    fontSize: typography.mood.fontSize,
    lineHeight: typography.mood.lineHeight,
    color: colors.primary,
    marginBottom: spacing.sm,
  },
  subtitle: {
    fontSize: typography.description.fontSize,
    lineHeight: typography.description.lineHeight,
    color: colors.secondary,
    textAlign: 'center',
  },
  actions: {
    gap: spacing.md,
    width: '100%',
  },
  permissionHint: {
    marginTop: spacing.lg,
    fontSize: typography.meta.fontSize,
    lineHeight: typography.meta.lineHeight,
    color: colors.secondary,
    textAlign: 'center',
  },
  permissionRow: {
    marginTop: spacing.lg,
    alignItems: 'center',
    gap: spacing.sm,
  },
});
