import * as SecureStore from 'expo-secure-store';

// Mirror of @clerk/expo's TokenCache interface.
interface TokenCache {
  getToken(key: string): Promise<string | null>;
  saveToken(key: string, token: string): Promise<void>;
}

export const clerkTokenCache: TokenCache = {
  async getToken(key: string) {
    try {
      return SecureStore.getItemAsync(key);
    } catch {
      return null;
    }
  },
  async saveToken(key: string, token: string) {
    try {
      return SecureStore.setItemAsync(key, token);
    } catch {
      // Ignore secure-store errors (e.g. device without hardware encryption).
    }
  },
};
