import 'package:dio/dio.dart';
import 'home_api.dart';

class HomeRepository {
  final HomeApi _api = HomeApi();

  Future<Map<String, dynamic>> getPosts({String? afterId}) async {
    return await _api.fetchPosts(afterId: afterId);
  }

  Future<Map<String, dynamic>> getPostDetail(String postId) async {
    return await _api.fetchPostDetail(postId);
  }

  Future<void> toggleLike(String postId) async {
    await _api.toggleLike(postId);
  }

  Future<Map<String, dynamic>> getComments(String postId) async =>
      await _api.fetchComments(postId);

  Future<Map<String, dynamic>> addComment(String postId, String text) async =>
      await _api.addComment(postId, text);

  Future<Map<String, dynamic>> updateComment(
    String commentId,
    String text,
  ) async => await _api.updateComment(commentId, text);

  Future<Map<String, dynamic>> deleteComment(String commentId) async =>
      await _api.deleteComment(commentId);

  Future<Map<String, dynamic>> createPost({
    required String content,
    required String visibility,
    String? documentType,
    List<MultipartFile>? images,
    MultipartFile? video,
    List<MultipartFile>? documents,
  }) async {
    return await _api.createPost(
      content: content,
      visibility: visibility,
      documentType: documentType,
      images: images,
      video: video,
      documents: documents,
    );
  }
}
