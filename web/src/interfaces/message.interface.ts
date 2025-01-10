export interface Message {
  id: string;
  channelId: string;
  authorId: string;
  content: string;
  /**
   * @kind iso8601
   */
  timestamp: string;
  read: boolean;
  attachments: MessageAttachment[];
}

export interface NewMessageAttachment {
  type: string;
  filename: string;
}

export interface MessageAttachment {
  id: string;
  messageId: string;
  filename: string;
  type: string;
  url?: string;
}
