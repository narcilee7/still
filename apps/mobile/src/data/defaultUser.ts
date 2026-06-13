import { User } from '@still/shared-types';

export interface CurrentUser extends User {
  postsCount: number;
  resonancesCount: number;
}

// Placeholder used only while the app is waiting for the authenticated session.
// The real user is loaded from the backend via /still.v1.UserService/GetMe.
export const GUEST_USER: CurrentUser = {
  id: '',
  username: 'guest',
  avatarUrl: undefined,
  postsCount: 0,
  resonancesCount: 0,
};
