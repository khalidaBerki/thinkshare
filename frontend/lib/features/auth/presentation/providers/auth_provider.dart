import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../data/auth_api.dart';
import 'package:go_router/go_router.dart';

class AuthProvider extends ChangeNotifier {
  final AuthApi _authApi = AuthApi();

  Future<void> login(
    String email,
    String password,
    BuildContext context,
  ) async {
    try {
      final token = await _authApi.login(email, password);
      debugPrint('Token: $token');

      // ✅ Save token
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('auth_token', token);

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
}
