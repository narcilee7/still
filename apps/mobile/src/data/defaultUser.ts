import { User } from '@still/shared-types';

export interface CurrentUser extends User {
  postsCount: number;
  resonancesCount: number;
}

// Default user until authentication is implemented.
// Must match the seeded demo user in apps/backend/internal/db/seed.go.
export const DEFAULT_USER: CurrentUser = {
  id: '00000000-0000-0000-0000-000000000001',
  username: 'still_moments',
  avatarUrl: undefined,
  postsCount: 0,
  resonancesCount: 0,
};
