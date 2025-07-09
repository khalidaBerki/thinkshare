import 'package:flutter/material.dart';
import 'package:dio/dio.dart';

import 'package:go_router/go_router.dart';
import '../widgets/media_carousel.dart';
import '../widgets/comment_section.dart';
import '../../data/home_repository.dart';

class PostDetailScreen extends StatefulWidget {
  final String postId;

  const PostDetailScreen({super.key, required this.postId});

  @override
  State<PostDetailScreen> createState() => _PostDetailScreenState();
}

class _PostDetailScreenState extends State<PostDetailScreen> {
  Map<String, dynamic>? post;
  bool isLoading = true;
  final HomeRepository _repository = HomeRepository();
  late bool hasLiked;
  late int likeCount;

  List<Map<String, dynamic>> comments = [];
  bool isCommentsLoading = false;

  @override
  void initState() {
    super.initState();
    _loadPost();
    _loadComments();
  }

  Future<void> _loadPost() async {
    try {
      final data = await _repository.getPostDetail(widget.postId);
      setState(() {
        post = data;
        isLoading = false;
        hasLiked = post?['user_has_liked'] == true;
        likeCount = post?['like_count'] ?? 0;
      });
    } catch (e) {
      debugPrint('Failed to load post detail: $e');
      setState(() {
        isLoading = false;
      });
    }
  }

  Future<void> _toggleLike() async {
    setState(() {
      hasLiked = !hasLiked;
      likeCount += hasLiked ? 1 : -1;
    });
    await _repository.toggleLike(widget.postId);
  }

  Future<void> _loadComments() async {
    setState(() => isCommentsLoading = true);
    final data = await _repository.getComments(widget.postId);
    setState(() {
      comments = List<Map<String, dynamic>>.from(data['comments']);
      isCommentsLoading = false;
    });
  }

  Future<void> _addComment(String text) async {
    try {
      final data = await _repository.addComment(widget.postId, text);
      final newComment = Map<String, dynamic>.from(data['comment']);
      setState(() {
        comments.insert(0, newComment);
      });
    } catch (e) {
      _showError(context, e.toString());
    }
  }

  Future<void> _editComment(String commentId, String text) async {
    try {
      final data = await _repository.updateComment(commentId, text);
      final updated = Map<String, dynamic>.from(data['comment']);
      setState(() {
        final idx = comments.indexWhere((c) => c['id'].toString() == commentId);
        if (idx != -1) comments[idx] = updated;
      });
    } catch (e) {
      _showError(context, _extractApiError(e));
    }
  }

  Future<void> _deleteComment(String commentId) async {
    try {
      await _repository.deleteComment(commentId);
      setState(() {
        comments.removeWhere((c) => c['id'].toString() == commentId);
      });
    } catch (e) {
      _showError(context, _extractApiError(e));
    }
  }

  void _showError(BuildContext context, String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  // Optionnel : pour extraire le message d'erreur de l'API (si Dio ou http)
  String _extractApiError(Object error) {
    if (error is DioException) {
      try {
        final data = error.response?.data;
        if (data is Map && data['error'] != null) {
          return data['error'].toString();
        }
      } catch (_) {}
      return error.message ?? 'Erreur inconnue';
    }
    // Sinon, retourne le message brut
    return error.toString();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    if (isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    if (post == null) {
      return const Scaffold(body: Center(child: Text("Post not found")));
    }

    final creator = post!['creator'] ?? {};
    final mediaUrls = List<String>.from(post!['media_urls'] ?? []);
    final isPrivate = post!['visibility'] == 'private';

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () {
            if (Navigator.of(context).canPop()) {
              Navigator.of(context).pop();
            } else {
              context.go('/home');
            }
          },
        ),
        title: const Text("Post Detail"),
        centerTitle: true,
      ),
      backgroundColor: colorScheme.surface,
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(18.0),
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 600),
            child: Material(
              elevation: 3,
              borderRadius: BorderRadius.circular(22),
              color: colorScheme.surface,
              child: Padding(
                padding: const EdgeInsets.all(18.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ListTile(
                      leading: InkWell(
                        borderRadius: BorderRadius.circular(28),
                        onTap: () {
                          if (creator['id'] != null) {
                            context.go('/user/${creator['id']}');
                          }
                        },
                        child: AnimatedScale(
                          scale: 1.0,
                          duration: const Duration(milliseconds: 100),
                          child: CircleAvatar(
                            backgroundImage: NetworkImage(
                              creator['avatar_url']?.isNotEmpty == true
                                  ? 'https://www.thinkshare.com/${creator['avatar_url']}'
                                  : 'https://ui-avatars.com/api/?name=${creator['full_name'] ?? 'User'}',
                            ),
                            radius: 28,
                          ),
                        ),
                      ),
                      title: InkWell(
                        borderRadius: BorderRadius.circular(4),
                        onTap: () {
                          if (creator['id'] != null) {
                            context.go('/user/${creator['id']}');
                          }
                        },
                        child: Padding(
                          padding: const EdgeInsets.symmetric(vertical: 2.0),
                          child: Text(
                            creator['full_name'] ?? 'No Name',
                            style: TextStyle(
                              fontFamily: 'Montserrat',
                              fontWeight: FontWeight.bold,
                              color: colorScheme.primary,
                              fontSize: 18,
                            ),
                          ),
                        ),
                      ),
                      subtitle: Row(
                        children: [
                          Text(
                            post!['document_type'] ?? '',
                            style: TextStyle(
                              fontFamily: 'Montserrat',
                              color: colorScheme.secondary,
                              fontSize: 13,
                            ),
                          ),
                          const SizedBox(width: 8),
                          Text(
                            _formatDate(post!['created_at']),
                            style: TextStyle(
                              color: colorScheme.secondary.withOpacity(0.7),
                              fontSize: 12,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 10),
                    if (isPrivate)
                      _PrivatePostBanner()
                    else ...[
                      Text(
                        post!['content'] ?? '',
                        style: TextStyle(
                          fontFamily: 'Montserrat',
                          fontSize: 16,
                          color: colorScheme.onSurface,
                        ),
                      ),
                      const SizedBox(height: 16),
                      if (mediaUrls.isNotEmpty) MediaCarousel(mediaUrls: mediaUrls),
                    ],
                    const SizedBox(height: 20),
                    Row(
                      children: [
                        IconButton(
                          icon: Icon(
                            hasLiked ? Icons.star : Icons.star_border,
                            color: colorScheme.primary,
                          ),
                          onPressed: _toggleLike,
                        ),
                        Text('$likeCount'),
                        const SizedBox(width: 16),
                        Icon(Icons.mode_comment_outlined, color: colorScheme.primary),
                        const SizedBox(width: 4),
                        Text('${post!['comment_count'] ?? 0}'),
                      ],
                    ),
                    const SizedBox(height: 20),
                    CommentSection(
                      comments: comments,
                      isLoading: isCommentsLoading,
                      onAddComment: _addComment,
                      onEditComment: _editComment,
                      onDeleteComment: _deleteComment,
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  String _formatDate(String? iso) {
    if (iso == null) return '';
    final date = DateTime.tryParse(iso);
    if (date == null) return '';
    return '${date.day}/${date.month}/${date.year}';
  }
}

class _PrivatePostBanner extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 16),
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(
        color: colorScheme.error.withOpacity(0.08),
        border: Border.all(color: colorScheme.error),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        children: [
          Text(
            "To view this private post, you must upgrade to premium",
            style: TextStyle(
              color: colorScheme.error,
              fontFamily: 'Montserrat',
              fontWeight: FontWeight.bold,
              fontSize: 15,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 10),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: colorScheme.error,
              foregroundColor: colorScheme.onError,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 10),
            ),
            onPressed: () {
              // TODO: Upgrade action
            },
            child: const Text("Upgrade to premium"),
          ),
        ],
      ),
    );
  }
}
