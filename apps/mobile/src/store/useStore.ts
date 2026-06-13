import { create } from 'zustand';
import { Post, User } from '@still/shared-types';
import { MOCK_USER, CurrentUser } from '../data/mockUser';

interface AppState {
  posts: Post[];
  user: CurrentUser;
  resonatedPostIds: Set<string>;

  setPosts: (posts: Post[]) => void;
  addPost: (post: Post) => void;
  updatePost: (postId: string, patch: Partial<Post>) => void;
  setUser: (user: CurrentUser) => void;
  toggleResonate: (postId: string) => void;
}

export const useStore = create<AppState>((set) => ({
  posts: [],
  user: MOCK_USER,
  resonatedPostIds: new Set(),

  setPosts: (posts) => set({ posts }),

  addPost: (post) =>
    set((state) => {
      const nextPosts = [post, ...state.posts];
      const userPosts = nextPosts.filter((p) => p.userId === state.user.id);
      return {
        posts: nextPosts,
        user: {
          ...state.user,
          postsCount: userPosts.length,
        },
      };
    }),

  updatePost: (postId, patch) =>
    set((state) => ({
      posts: state.posts.map((p) => (p.id === postId ? { ...p, ...patch } : p)),
    })),

  setUser: (user) => set({ user }),

  toggleResonate: (postId) =>
    set((state) => {
      const nextResonated = new Set(state.resonatedPostIds);
      const isResonating = !nextResonated.has(postId);

      const nextPosts = state.posts.map((post) => {
        if (post.id !== postId) return post;
        return {
          ...post,
          resonanceCount: isResonating
            ? post.resonanceCount + 1
            : Math.max(0, post.resonanceCount - 1),
        };
      });

      if (isResonating) {
        nextResonated.add(postId);
      } else {
        nextResonated.delete(postId);
      }

      const delta = isResonating ? 1 : -1;
      return {
        posts: nextPosts,
        resonatedPostIds: nextResonated,
        user: {
          ...state.user,
          resonancesCount: Math.max(0, state.user.resonancesCount + delta),
        },
      };
    }),
}));
