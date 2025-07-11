import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:go_router/go_router.dart';
import '../../data/profile_repository.dart';
import '../../data/profile_api.dart';

class ProfileProvider extends ChangeNotifier {
  final ProfileRepository _repo = ProfileRepository(ProfileApi());

  Map<String, dynamic>? myProfile;
  Map<String, dynamic>? userProfile;
  List<Map<String, dynamic>> myPosts = [];
  List<Map<String, dynamic>> userPosts = [];
  bool isLoading = false;
  bool isPostsLoading = false;
  String? error;
  String? postsError;

  bool isUpdating = false;
  String? updateError;
  String? updateSuccess;

  Future<void> fetchMyProfile() async {
    isLoading = true;
    error = null;
    notifyListeners();
    try {
      myProfile = await _repo.getMyProfile();
    } catch (e) {
      error = e.toString();
      myProfile = null;
    }
    isLoading = false;
    notifyListeners();
  }

  Future<void> fetchUserProfile(int userId) async {
    isLoading = true;
    error = null;
    notifyListeners();
    try {
      userProfile = await _repo.getUserProfile(userId);
    } catch (e) {
      error = e.toString();
      userProfile = null;
    }
    isLoading = false;
    notifyListeners();
  }

  Future<void> fetchMyPosts({int? after, int? limit}) async {
    isPostsLoading = true;
    postsError = null;
    notifyListeners();
    try {
      final res = await _repo.getUserPosts(
        myProfile?['id'] ?? 0,
        after: after,
        limit: limit,
      );
      myPosts = List<Map<String, dynamic>>.from(res['posts'] ?? []);
    } catch (e) {
      postsError = e.toString();
      myPosts = [];
    }
    isPostsLoading = false;
    notifyListeners();
  }

  Future<void> fetchUserPosts(int userId, {int? after, int? limit}) async {
    isPostsLoading = true;
    postsError = null;
    notifyListeners();
    try {
      final res = await _repo.getUserPosts(userId, after: after, limit: limit);
      userPosts = List<Map<String, dynamic>>.from(res['posts'] ?? []);
    } catch (e) {
      postsError = e.toString();
      userPosts = [];
    }
    isPostsLoading = false;
    notifyListeners();
  }

  Future<void> updateProfile({
    String? fullName,
    String? bio,
    String? avatarUrl,
  }) async {
    isUpdating = true;
    updateError = null;
    updateSuccess = null;
    notifyListeners();
    try {
      final data = <String, dynamic>{};
      if (fullName != null) data['full_name'] = fullName;
      if (bio != null) data['bio'] = bio;
      if (avatarUrl != null) data['avatar_url'] = avatarUrl;
      if (data.isEmpty) throw Exception("Aucun champ à modifier");

      final res = await _repo.updateProfile(data);
      updateSuccess = res['message'] ?? "Profil mis à jour";
      await fetchMyProfile();
    } catch (e) {
      updateError = e.toString();
    }
    isUpdating = false;
    notifyListeners();
  }

  Future<void> logout(BuildContext context) async {
    try {
      await _repo.logout();
    } catch (_) {}
    // Nettoie les prefs si besoin
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('auth_token');
    await prefs.remove('user_id');
    // Redirige vers la page d'accueil
    if (context.mounted) context.go('/');
  }
}
