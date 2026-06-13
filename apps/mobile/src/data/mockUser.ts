import { User } from '@still/shared-types';

export interface CurrentUser extends User {
  postsCount: number;
  resonancesCount: number;
}

export const MOCK_USER: CurrentUser = {
  id: 'user-1',
  username: 'still_moments',
  avatarUrl: undefined,
  postsCount: 3,
  resonancesCount: 47,
};
