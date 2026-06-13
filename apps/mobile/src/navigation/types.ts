export type MainTabParamList = {
  Feed: undefined;
  Create: undefined;
  Profile: undefined;
};

export type CreateStackParamList = {
  CreateSelect: undefined;
  CreateEdit: { imageUri: string };
  CreateSuccess: { postId: string };
};

export type RootStackParamList = {
  Main: undefined;
  PostDetail: { postId: string };
  EditPost: { postId: string };
};
