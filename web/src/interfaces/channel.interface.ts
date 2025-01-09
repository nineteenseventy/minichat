interface BaseChannel {
  id: string;
  /**
   * @kind iso8601
   */
  createdAt: string;
  unreadCount: number;
}

export interface DirectChannel extends BaseChannel {
  type: 'direct';
  title: string;
}

export interface PublicChannel extends BaseChannel {
  type: 'public';
  title: string;
  description: string;
}

export interface GroupChannel extends BaseChannel {
  type: 'group';
  title: string;
}

export type Channel = DirectChannel | PublicChannel | GroupChannel;

export interface ChannelsResponse {
  private: (GroupChannel | DirectChannel)[];
  public: PublicChannel[];
}
