import '../data/message_api.dart';

class MessageRepository {
  final MessageApi api;
  MessageRepository(this.api);

  Future<List<Map<String, dynamic>>> getConversations() =>
      api.getConversations();
  Future<List<Map<String, dynamic>>> getMessagesWithUser(int otherUserId) =>
      api.getMessagesWithUser(otherUserId);
  Future<Map<String, dynamic>> sendMessage(int receiverId, String content) =>
      api.sendMessage(receiverId, content);
  Future<void> markAsRead(int senderId) => api.markAsRead(senderId);
  Future<Map<String, dynamic>> updateMessage(int messageId, String content) =>
      api.updateMessage(messageId, content);
  Future<void> deleteMessage(int messageId) => api.deleteMessage(messageId);
}
