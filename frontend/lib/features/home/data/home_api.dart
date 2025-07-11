import '../../../services/api_service.dart';
import 'package:dio/dio.dart';

class HomeApi {
  final dio = ApiService().dio;

  Future<Map<String, dynamic>> fetchPosts({String? afterId}) async {
    final params = afterId != null ? {'after': afterId} : null;
    final response = await dio.get('api/posts', queryParameters: params);
    return response.data;
  }

  Future<Map<String, dynamic>> fetchPostDetail(String postId) async {
    final response = await dio.get('api/posts/$postId');
    return response.data;
  }

  Future<void> toggleLike(String postId) async {
    await dio.post('api/likes/posts/$postId');
  }

  // GET /api/comments/:postID
  Future<Map<String, dynamic>> fetchComments(String postId) async {
    final response = await dio.get('api/comments/$postId');
    return response.data;
  }

  // POST /api/comments
  Future<Map<String, dynamic>> addComment(String postId, String text) async {
    final response = await dio.post(
      'api/comments',
      data: {'post_id': int.tryParse(postId) ?? postId, 'text': text},
    );
    return response.data;
  }

  // PUT /api/comments/:id
  Future<Map<String, dynamic>> updateComment(
    String commentId,
    String text,
  ) async {
    final response = await dio.put(
      'api/comments/$commentId',
      data: {'text': text},
    );
    return response.data;
  }

  // DELETE /api/comments/:id
  Future<Map<String, dynamic>> deleteComment(String commentId) async {
    final response = await dio.delete('api/comments/$commentId');
    return response.data;
  }

  // POST /api/posts
  Future<Map<String, dynamic>> createPost({
    required String content,
    required String visibility,
    String? documentType,
    bool isPaidOnly = false,
    List<MultipartFile>? images,
    MultipartFile? video,
    List<MultipartFile>? documents,
  }) async {
    final formData = FormData();

    formData.fields
      ..add(MapEntry('content', content))
      ..add(MapEntry('visibility', visibility))
      ..add(MapEntry('is_paid_only', isPaidOnly.toString()));
    if (documentType != null) {
      formData.fields.add(MapEntry('document_type', documentType));
    }
    if (images != null && images.isNotEmpty) {
      for (final img in images) {
        formData.files.add(MapEntry('images', img));
      }
    }
    if (video != null) {
      formData.files.add(MapEntry('video', video));
    }
    if (documents != null && documents.isNotEmpty) {
      for (final doc in documents) {
        formData.files.add(MapEntry('documents', doc));
      }
    }

    final response = await dio.post(
      'api/posts/',
      data: formData,
      options: Options(contentType: 'multipart/form-data'),
    );
    return response.data;
  }
}
