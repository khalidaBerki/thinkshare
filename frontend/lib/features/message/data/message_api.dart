import '../../../services/api_service.dart';

class MessageApi {
  final dio = ApiService().dio;

  // Liste des conversations
  Future<List<Map<String, dynamic>>> getConversations() async {
    final response = await dio.get('api/messages/conversations');
    return List<Map<String, dynamic>>.from(response.data);
  }

  // Récupérer les messages d'une conversation
  Future<List<Map<String, dynamic>>> getMessagesWithUser(
    int otherUserId,
  ) async {
    final response = await dio.get('api/messages/$otherUserId');
    return List<Map<String, dynamic>>.from(response.data);
  }

  // Envoyer un message
  Future<Map<String, dynamic>> sendMessage(
    int receiverId,
    String content,
  ) async {
    final response = await dio.post(
      'api/messages',
      data: {'receiver_id': receiverId, 'content': content},
    );
    return response.data;
  }

  // Marquer comme lu
  Future<void> markAsRead(int senderId) async {
    await dio.patch('api/messages/$senderId/read');
  }

  // Modifier un message
  Future<Map<String, dynamic>> updateMessage(
    int messageId,
    String content,
  ) async {
    final response = await dio.put(
      'api/messages/$messageId',
      data: {'content': content},
    );
    return response.data;
  }

  // Supprimer un message
  Future<void> deleteMessage(int messageId) async {
    await dio.delete('api/messages/$messageId');
  }
}
