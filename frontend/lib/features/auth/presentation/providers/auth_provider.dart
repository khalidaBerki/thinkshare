import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../data/auth_api.dart';
import 'package:go_router/go_router.dart';

class AuthProvider extends ChangeNotifier {
  final AuthApi _authApi = AuthApi();

  int? _userId;
  int? get userId => _userId;

  Future<void> login(
    String email,
    String password,
    BuildContext context,
  ) async {
    try {
      // On suppose que _authApi.login renvoie { "token": "...", "user_id": 42 }
      final response = await _authApi.login(email, password);
      debugPrint('Login response: $response');

      final token = response['token'];
      final userIdRaw = response['user_id'];

      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('auth_token', token);

      int? userId;
      if (userIdRaw is int) {
        userId = userIdRaw;
      } else if (userIdRaw is String) {
        userId = int.tryParse(userIdRaw);
      }

      if (userId != null) {
        _userId = userId;
        await prefs.setInt('user_id', _userId!);
      }

      notifyListeners();

      if (!context.mounted) return;
      context.go('/home');
    } catch (e) {
      debugPrint('Login failed: $e');
    }
  }

  Future<void> register({
    required String name,
    required String firstname,
    required String username,
    required String email,
    required String password,
    required BuildContext context,
  }) async {
    try {
      await _authApi.register(
        name: name,
        firstname: firstname,
        username: username,
        email: email,
        password: password,
      );

      if (!context.mounted) return;
      context.go('/login');
    } catch (e) {
      debugPrint('Register failed: $e');
    }
  }

  Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('auth_token');
  }

  Future<void> loadUserId() async {
    final prefs = await SharedPreferences.getInstance();
    _userId = prefs.getInt('user_id');
    notifyListeners();
  }
}
