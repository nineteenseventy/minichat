export interface UserProfile {
  id: string;
  username: string;
  bio?: string;
  picture?: string;
}

export interface UpdateUserProfilePayload {
  username?: string;
  bio?: string;
}
