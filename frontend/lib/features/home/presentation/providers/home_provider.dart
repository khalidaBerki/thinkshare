import 'package:flutter/material.dart';
import '../../data/home_repository.dart';

class HomeProvider extends ChangeNotifier {
  final HomeRepository _repository = HomeRepository();

  List<Map<String, dynamic>> posts = [];
  bool isLoading = false;
  bool hasMore = true;
  String? lastId;

  Future<void> loadPosts({bool refresh = false}) async {
    if (isLoading) return;
    isLoading = true;
    notifyListeners();

    try {
      if (refresh) {
        posts.clear();
        lastId = null;
        hasMore = true;
      }

      final data = await _repository.getPosts(afterId: lastId);
      final newPosts = List<Map<String, dynamic>>.from(data['posts']);

      posts.addAll(newPosts);
      lastId = data['last_id']?.toString();
      hasMore = data['has_more'] ?? false;
    } catch (e) {
      debugPrint('Failed to load posts: $e');
    }

    isLoading = false;
    notifyListeners();
  }
}
