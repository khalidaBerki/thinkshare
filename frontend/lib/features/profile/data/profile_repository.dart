import 'profile_api.dart';

class ProfileRepository {
  final ProfileApi api;
  ProfileRepository(this.api);

  Future<Map<String, dynamic>> getMyProfile() => api.getMyProfile();
  Future<Map<String, dynamic>> getUserProfile(int userId) =>
      api.getUserProfile(userId);
  Future<Map<String, dynamic>> getUserPosts(int userId, {int? after, int? limit}) =>
      api.getUserPosts(userId, after: after, limit: limit);
  Future<Map<String, dynamic>> updateProfile(Map<String, dynamic> data) => api.updateProfile(data);
  Future<void> logout() => api.logout();
}
