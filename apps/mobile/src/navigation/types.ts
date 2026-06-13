export type RootTabParamList = {
  Feed: undefined;
  Create: undefined;
  Profile: undefined;
};

export type CreateStackParamList = {
  CreateSelect: undefined;
  CreateEdit: { imageUri: string };
  CreateSuccess: { postId: string };
};
