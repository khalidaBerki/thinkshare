import 'package:dio/dio.dart';
import '../../../config/api_config.dart';

class AuthApi {
  final Dio _dio = Dio(BaseOptions(baseUrl: ApiConfig.baseUrl));

  Future<String> login(String email, String password) async {
    final response = await _dio.post(
      '/login',
      data: {'email': email, 'password': password},
    );
    return response.data['token'];
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
