import 'package:dio/dio.dart';
import '../../../config/api_config.dart';

class AuthApi {
  final Dio _dio = Dio(BaseOptions(baseUrl: ApiConfig.baseUrl));

  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await _dio.post(
      '/login',
      data: {'email': email, 'password': password},
    );
    // On retourne tout le JSON (token + user_id)
    return Map<String, dynamic>.from(response.data);
  }

  Future<void> register({
    required String name,
    required String firstname,
    required String username,
    required String email,
    required String password,
  }) async {
    await _dio.post(
      '/register',
      data: {
        'name': name,
        'firstname': firstname,
        'username': username,
        'email': email,
        'password': password,
      },
    );
  }
}
