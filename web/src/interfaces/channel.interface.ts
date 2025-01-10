export interface Channel {
  id: string;
  /**
   * @kind iso8601
   */
  createdAt: string;
  unreadCount: number;
  type: 'direct' | 'public' | 'group';
  title: string;
  description?: string;
}
