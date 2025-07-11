import '../../home/data/home_api.dart';
import '../../message/data/message_api.dart';

class AdminDashboardRepository {
  final HomeApi _homeApi = HomeApi();
  final MessageApi _messageApi = MessageApi();

  // Récupère tous les posts
  Future<List<Map<String, dynamic>>> getAllPosts() async {
    final response = await _homeApi.fetchPosts();
    return List<Map<String, dynamic>>.from(response['posts']);
  }

  // Récupère tous les commentaires d'un post
  Future<List<Map<String, dynamic>>> getCommentsForPost(String postId) async {
    final response = await _homeApi.fetchComments(postId);
    return List<Map<String, dynamic>>.from(response['comments']);
  }

  // Récupère toutes les conversations/messages
  Future<List<Map<String, dynamic>>> getAllConversations() async {
    return await _messageApi.getConversations();
  }

  // Supprimer un commentaire (si tu ajoutes la méthode dans HomeApi)
  Future<void> deleteComment(String commentId) async {
    await _homeApi.deleteComment(commentId);
  }
}
