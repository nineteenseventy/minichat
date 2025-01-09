export interface User {
  id: string;
  username: string;
  picture?: string;
}

export type UserStatus = 'online' | 'offline' | 'away';
export interface UserStatusResponse {
  id: string;
  status: UserStatus;
}
