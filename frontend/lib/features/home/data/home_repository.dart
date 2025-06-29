import 'home_api.dart';

class HomeRepository {
  final HomeApi _api = HomeApi();

  Future<Map<String, dynamic>> getPosts({String? afterId}) async {
    return await _api.fetchPosts(afterId: afterId);
  }

  Future<Map<String, dynamic>> getPostDetail(String postId) async {
    return await _api.fetchPostDetail(postId);
  }
}
