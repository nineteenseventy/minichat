import type { User } from './user.interface';

export interface Message {
  id: string;
  author: User;
  content: string;
  timestamp: string;
}
