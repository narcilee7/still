export const MOODS = [
  'still',
  'waiting',
  'drift',
  'warm',
  'distant',
  'quiet',
  'hollow',
  'soft',
  'fading',
  'wondering',
  'returning',
] as const;

export type Mood = (typeof MOODS)[number];

export interface Post {
  id: string;
  userId: string;
  imageUrl: string;
  mood: Mood;
  title: string;
  description: string;
  createdAt: string;
  resonanceCount: number;
}

export interface User {
  id: string;
  username: string;
  avatarUrl?: string;
}
