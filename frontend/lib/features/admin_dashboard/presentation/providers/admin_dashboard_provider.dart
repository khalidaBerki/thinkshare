import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../data/admin_dashboard_repository.dart';

class AdminDashboardProvider extends ChangeNotifier {
  final AdminDashboardRepository _repo = AdminDashboardRepository();

  bool isLoading = false;
  DateTimeRange? selectedRange;

  // Stats
  int totalPosts = 0;
  int totalComments = 0;
  int totalConversations = 0;
  Map<String, int> postsPerDay = {};
  List<Map<String, dynamic>> topPosts = [];
  List<Map<String, dynamic>> flopPosts = [];

  // Ajoute cette propriété
  Map<String, Map<String, int>> mediaProgression = {};

  Future<void> fetchStats() async {
    isLoading = true;
    notifyListeners();

    // Récupère toutes les données
    final posts = await _repo.getAllPosts();
    final conversations = await _repo.getAllConversations();

    // Filtrage par date si besoin
    List<Map<String, dynamic>> filteredPosts = posts;
    if (selectedRange != null) {
      filteredPosts = posts.where((post) {
        final createdAt = DateTime.tryParse(post['created_at'] ?? '');
        if (createdAt == null) return false;
        return createdAt.isAfter(
              selectedRange!.start.subtract(const Duration(days: 1)),
            ) &&
            createdAt.isBefore(selectedRange!.end.add(const Duration(days: 1)));
      }).toList();
    }

    // Calculs de base
    totalPosts = filteredPosts.length;
    totalComments = 0;
    for (final post in filteredPosts) {
      totalComments += ((post['comment_count'] ?? 0) as num).toInt();
    }
    totalConversations = conversations.length;

    // AJOUTE CE CODE ici pour media progression
    mediaProgression = {};
    for (final p in filteredPosts) {
      final date = DateTime.tryParse((p['created_at'] ?? '').toString());
      if (date != null) {
        final key = DateFormat('yyyy-MM-dd').format(date);
        mediaProgression.putIfAbsent(
          key,
          () => {'image': 0, 'video': 0, 'doc': 0},
        );

        final mediaUrls = p['media_urls'];
        if (mediaUrls is List && mediaUrls.isNotEmpty) {
          for (var url in mediaUrls) {
            if (url == null) continue;
            final urlStr = url.toString();
            if (urlStr.endsWith('.mp4') || urlStr.endsWith('.mov')) {
              mediaProgression[key]!['video'] =
                  (mediaProgression[key]!['video'] ?? 0) + 1;
            } else if (urlStr.endsWith('.pdf') || urlStr.endsWith('.doc')) {
              mediaProgression[key]!['doc'] =
                  (mediaProgression[key]!['doc'] ?? 0) + 1;
            } else {
              mediaProgression[key]!['image'] =
                  (mediaProgression[key]!['image'] ?? 0) + 1;
            }
          }
        }
      }
    }

    // Posts per day
    postsPerDay = {};
    for (final p in filteredPosts) {
      final date = DateTime.tryParse(p['created_at'] ?? '');
      if (date != null) {
        final key =
            "${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}";
        postsPerDay[key] = (postsPerDay[key] ?? 0) + 1;
      }
    }

    // Top/Flop posts
    topPosts = List<Map<String, dynamic>>.from(filteredPosts)
      ..sort(
        (a, b) => ((b['like_count'] ?? 0) + (b['comment_count'] ?? 0))
            .compareTo((a['like_count'] ?? 0) + (a['comment_count'] ?? 0)),
      );
    topPosts = topPosts.take(5).toList();

    flopPosts = List<Map<String, dynamic>>.from(filteredPosts)
      ..sort(
        (a, b) => ((a['like_count'] ?? 0) + (a['comment_count'] ?? 0))
            .compareTo((b['like_count'] ?? 0) + (b['comment_count'] ?? 0)),
      );
    flopPosts = flopPosts.take(5).toList();

    isLoading = false;
    notifyListeners();
  }

  void setRange(DateTimeRange? range) {
    selectedRange = range;
    fetchStats();
  }

  // Actions admin
  Future<void> deleteComment(String commentId) async {
    await _repo.deleteComment(commentId);
    await fetchStats();
  }
}
