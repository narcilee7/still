import { useStore } from './useStore';
import { DEFAULT_USER } from '../data/defaultUser';

const createPost = (id: string, overrides = {}) => ({
  id,
  userId: DEFAULT_USER.id,
  imageUrl: 'https://example.com/photo.jpg',
  mood: 'still' as const,
  title: 'A moment',
  description: 'A quiet description',
  createdAt: new Date().toISOString(),
  resonanceCount: 0,
  ...overrides,
});

describe('useStore', () => {
  beforeEach(() => {
    useStore.setState({
      posts: [],
      user: DEFAULT_USER,
      resonatedPostIds: new Set(),
    });
  });

  it('sets posts', () => {
    const posts = [createPost('1'), createPost('2')];
    useStore.getState().setPosts(posts);
    expect(useStore.getState().posts).toHaveLength(2);
  });

  it('adds a post to the top and updates postsCount', () => {
    useStore.getState().addPost(createPost('1'));
    const state = useStore.getState();
    expect(state.posts[0].id).toBe('1');
    expect(state.user.postsCount).toBe(1);
  });

  it('appends new posts without duplicates', () => {
    useStore.getState().setPosts([createPost('1')]);
    useStore.getState().appendPosts([createPost('1'), createPost('2')]);
    expect(useStore.getState().posts).toHaveLength(2);
  });

  it('updates a post', () => {
    useStore.getState().setPosts([createPost('1', { title: 'Old' })]);
    useStore.getState().updatePost('1', { title: 'New' });
    expect(useStore.getState().posts[0].title).toBe('New');
  });

  it('toggles resonance optimistically', () => {
    useStore.getState().setPosts([createPost('1')]);
    useStore.getState().toggleResonate('1');
    const state = useStore.getState();
    expect(state.resonatedPostIds.has('1')).toBe(true);
    expect(state.posts[0].resonanceCount).toBe(1);
    expect(state.user.resonancesCount).toBe(1);
  });

  it('toggles resonance off', () => {
    useStore.getState().setPosts([createPost('1')]);
    useStore.getState().toggleResonate('1');
    useStore.getState().toggleResonate('1');
    const state = useStore.getState();
    expect(state.resonatedPostIds.has('1')).toBe(false);
    expect(state.posts[0].resonanceCount).toBe(0);
    expect(state.user.resonancesCount).toBe(0);
  });

  it('setResonated ignores no-op changes', () => {
    useStore.getState().setPosts([createPost('1')]);
    useStore.getState().setResonated('1', false);
    expect(useStore.getState().posts[0].resonanceCount).toBe(0);
  });

  it('setResonated updates count to match server state', () => {
    useStore.getState().setPosts([createPost('1', { resonanceCount: 5 })]);
    useStore.getState().setResonated('1', true);
    const state = useStore.getState();
    expect(state.resonatedPostIds.has('1')).toBe(true);
    expect(state.posts[0].resonanceCount).toBe(6);
    expect(state.user.resonancesCount).toBe(1);
  });
});
