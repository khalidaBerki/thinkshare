import '../../../services/api_service.dart';

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
}
