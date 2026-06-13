import { Interceptor, createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import {
  AnalyzeService,
  CreatePostRequest,
  DeletePostRequest,
  FeedService,
  GetProfileRequest,
  GetUploadURLRequest,
  PostService,
  ResonateService,
  StorageService,
  UpdatePostRequest,
  UserService,
} from '@still/generated-sdk';
import { Mood, Post, User } from '@still/shared-types';
import { getToken } from './authToken';

const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

const authInterceptor: Interceptor = (next) => async (req) => {
  const token = await getToken();
  if (token) {
    req.header.set('Authorization', `Bearer ${token}`);
  }
  return next(req);
};

const transport = createConnectTransport({
  baseUrl: API_BASE_URL,
  interceptors: [authInterceptor],
});

const feedClient = createPromiseClient(FeedService, transport);
const postClient = createPromiseClient(PostService, transport);
const analyzeClient = createPromiseClient(AnalyzeService, transport);
const resonateClient = createPromiseClient(ResonateService, transport);
const userClient = createPromiseClient(UserService, transport);
const storageClient = createPromiseClient(StorageService, transport);

export interface AnalysisResult {
  mood: Mood;
  title: string;
  description: string;
}

export interface ProfileResult {
  user: User;
  postsCount: number;
  resonancesCount: number;
}

function mapProtoPost(p: unknown): Post {
  const proto = p as Record<string, unknown>;
  return {
    id: String(proto.id ?? ''),
    userId: String(proto.userId ?? ''),
    imageUrl: String(proto.imageUrl ?? ''),
    mood: String(proto.mood ?? 'still') as Mood,
    title: String(proto.title ?? ''),
    description: String(proto.description ?? ''),
    createdAt:
      (proto.createdAt as { toDate?: () => Date } | undefined)?.toDate?.().toISOString() ??
      new Date().toISOString(),
    resonanceCount: Number(proto.resonanceCount ?? 0),
  };
}

export interface FeedPage {
  posts: Post[];
  nextPageToken: string;
}

export async function listFeed(pageToken = ''): Promise<FeedPage> {
  const res = await feedClient.listFeed({ pageSize: 20, pageToken });
  return {
    posts: (res.posts ?? []).map(mapProtoPost),
    nextPageToken: res.nextPageToken ?? '',
  };
}

export async function analyzeImage(imageUrl: string): Promise<AnalysisResult> {
  const res = await analyzeClient.analyzeImage({ imageUrl });
  return {
    mood: (res.mood || 'still') as Mood,
    title: res.title || '',
    description: res.description || '',
  };
}

export async function createPost(payload: {
  imageUrl: string;
  mood: Mood;
  title: string;
  description: string;
}): Promise<Post> {
  const req = new CreatePostRequest({
    imageUrl: payload.imageUrl,
    mood: payload.mood,
    title: payload.title,
    description: payload.description,
  });
  const res = await postClient.createPost(req);
  return mapProtoPost(res.post);
}

export async function getPost(postId: string): Promise<Post> {
  const res = await postClient.getPost({ id: postId });
  return mapProtoPost(res.post);
}

export async function updatePost(payload: {
  id: string;
  mood: Mood;
  title: string;
  description: string;
}): Promise<Post> {
  const req = new UpdatePostRequest({
    id: payload.id,
    mood: payload.mood,
    title: payload.title,
    description: payload.description,
  });
  const res = await postClient.updatePost(req);
  return mapProtoPost(res.post);
}

export async function deletePost(postId: string): Promise<void> {
  await postClient.deletePost(new DeletePostRequest({ id: postId }));
}

export interface ResonateResult {
  post: Post;
  hasResonated: boolean;
}

export async function resonate(postId: string): Promise<ResonateResult> {
  const res = await resonateClient.resonate({ postId });
  return {
    post: mapProtoPost(res.post),
    hasResonated: res.hasResonated,
  };
}

export async function getMe(): Promise<User> {
  const res = await userClient.getMe({});
  return {
    id: res.user?.id || '',
    username: res.user?.username || '',
    avatarUrl: res.user?.avatarUrl || undefined,
  };
}

export async function getProfile(userId: string): Promise<ProfileResult> {
  const req = new GetProfileRequest({ userId });
  const res = await userClient.getProfile(req);
  const user: User = {
    id: res.user?.id || userId,
    username: res.user?.username || '',
    avatarUrl: res.user?.avatarUrl || undefined,
  };
  return {
    user,
    postsCount: res.postsCount ?? 0,
    resonancesCount: res.resonancesCount ?? 0,
  };
}

export async function getUploadURL(filename: string, contentType: string) {
  const req = new GetUploadURLRequest({ filename, contentType });
  const res = await storageClient.getUploadURL(req);
  return {
    uploadUrl: res.uploadUrl,
    publicUrl: res.publicUrl,
  };
}

export async function uploadImage(localUri: string, uploadUrl: string, contentType: string) {
  const blob = await fetch(localUri).then((r) => r.blob());
  const response = await fetch(uploadUrl, {
    method: 'PUT',
    body: blob,
    headers: {
      'Content-Type': contentType,
      'x-amz-acl': 'public-read',
    },
  });
  if (!response.ok) {
    const body = await response.text();
    throw new Error(`upload failed: ${response.status} ${body}`);
  }
}
