import { Mood, Post } from '@still/shared-types';
import { MOCK_POSTS } from '../data/mockPosts';
import { MOCK_USER } from '../data/mockUser';

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export interface AnalysisResult {
  mood: Mood;
  title: string;
  description: string;
}

export interface ProfileResult {
  user: {
    id: string;
    username: string;
    avatarUrl?: string;
  };
  posts: Post[];
  postsCount: number;
  resonancesCount: number;
}

export async function listFeed(): Promise<Post[]> {
  await delay(400);
  return [...MOCK_POSTS];
}

export async function analyzeImage(): Promise<AnalysisResult> {
  await delay(1200);
  return {
    mood: 'waiting',
    title: 'On The Way',
    description: 'Some days are made of unfinished thoughts.',
  };
}

export async function createPost(payload: {
  imageUrl: string;
  mood: Mood;
  title: string;
  description: string;
}): Promise<Post> {
  await delay(800);
  return {
    id: `post-${Date.now()}`,
    userId: MOCK_USER.id,
    imageUrl: payload.imageUrl,
    mood: payload.mood,
    title: payload.title,
    description: payload.description,
    createdAt: new Date().toISOString(),
    resonanceCount: 0,
  };
}

export async function getProfile(): Promise<ProfileResult> {
  await delay(300);
  const userPosts = MOCK_POSTS.filter((p) => p.userId === MOCK_USER.id);
  return {
    user: {
      id: MOCK_USER.id,
      username: MOCK_USER.username,
      avatarUrl: MOCK_USER.avatarUrl,
    },
    posts: userPosts,
    postsCount: userPosts.length,
    resonancesCount: MOCK_USER.resonancesCount,
  };
}
