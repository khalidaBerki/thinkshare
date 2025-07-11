import '../../../services/api_service.dart';

class ProfileApi {
  final dio = ApiService().dio;

  Future<Map<String, dynamic>> updateProfile(Map<String, dynamic> data) async {
    final response = await dio.put('/api/profile', data: data);
    return response.data;
  }

  Future<void> logout() async {
    await dio.get('/logout');
  }

  Future<Map<String, dynamic>> getMyProfile() async {
    final response = await dio.get('/api/profile');
    return response.data;
  }

  Future<Map<String, dynamic>> getUserProfile(int userId) async {
    final response = await dio.get('/api/users/$userId');
    return response.data;
  }

  Future<Map<String, dynamic>> getUserPosts(
    int userId, {
    int? after,
    int? limit,
  }) async {
    final query = <String, dynamic>{};
    if (after != null) query['after'] = after;
    if (limit != null) query['limit'] = limit;
    final response = await dio.get(
      '/api/posts/user/$userId',
      queryParameters: query,
    );
    return response.data;
  }
}
