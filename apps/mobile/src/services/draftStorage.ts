import AsyncStorage from '@react-native-async-storage/async-storage';
import { Mood } from '@still/shared-types';

const DRAFT_KEY = '@still/createDraft';

export interface CreateDraft {
  imageUri: string;
  publicUrl?: string;
  mood: Mood;
  title: string;
  description: string;
  savedAt: string;
}

export async function loadDraft(): Promise<CreateDraft | null> {
  try {
    const raw = await AsyncStorage.getItem(DRAFT_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as CreateDraft;
  } catch {
    return null;
  }
}

export async function saveDraft(draft: CreateDraft): Promise<void> {
  try {
    await AsyncStorage.setItem(DRAFT_KEY, JSON.stringify(draft));
  } catch {
    // Ignore storage errors.
  }
}

export async function clearDraft(): Promise<void> {
  try {
    await AsyncStorage.removeItem(DRAFT_KEY);
  } catch {
    // Ignore storage errors.
  }
}
