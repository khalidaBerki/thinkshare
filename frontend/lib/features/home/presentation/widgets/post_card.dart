import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/services.dart'; // Pour Clipboard
import 'media_carousel.dart';

class PostCard extends StatelessWidget {
  final Map<String, dynamic> post;

  const PostCard({super.key, required this.post});

  @override
  Widget build(BuildContext context) {
    final creator = post['creator'] ?? {};
    final mediaUrls = List<String>.from(post['media_urls'] ?? []);
    final isPrivate = post['visibility'] == 'private';
    final colorScheme = Theme.of(context).colorScheme;
    final postId = post['id'].toString();

    return InkWell(
      borderRadius: BorderRadius.circular(18),
      onTap: () {
        context.go('/post/$postId');
      },
      child: Card(
        margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        elevation: 3,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ListTile(
                leading: InkWell(
                  borderRadius: BorderRadius.circular(24),
                  onTap: () {
                    if (creator['id'] != null) {
                      context.go('/user/${creator['id']}');
                    }
                  },
                  child: CircleAvatar(
                    backgroundImage: NetworkImage(
                      creator['avatar_url']?.isNotEmpty == true
                          ? 'https://www.thinkshare.com/${creator['avatar_url']}'
                          : 'https://ui-avatars.com/api/?name=${creator['full_name'] ?? 'User'}',
                    ),
                    radius: 24,
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
                        fontSize: 16,
                      ),
                    ),
                  ),
                ),
                subtitle: Row(
                  children: [
                    Text(
                      post['document_type'] ?? '',
                      style: TextStyle(
                        fontFamily: 'Montserrat',
                        color: colorScheme.secondary,
                        fontSize: 13,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      _formatDate(post['created_at']),
                      style: TextStyle(
                        color: colorScheme.secondary.withOpacity(0.7),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
                // onTap retiré ici pour laisser InkWell global gérer le tap
              ),
              if (isPrivate)
                _PrivatePostBanner()
              else ...[
                Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16.0,
                    vertical: 4,
                  ),
                  child: Text(
                    post['content'] ?? '',
                    style: TextStyle(
                      fontFamily: 'Montserrat',
                      fontSize: 15,
                      color: colorScheme.onSurface,
                    ),
                  ),
                ),
                if (mediaUrls.isNotEmpty) MediaCarousel(mediaUrls: mediaUrls),
              ],
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 16.0,
                  vertical: 8,
                ),
                child: Row(
                  children: [
                    IconButton(
                      icon: Icon(
                        post['user_has_liked'] == true
                            ? Icons.star
                            : Icons.star_border,
                        color: colorScheme.primary,
                      ),
                      onPressed: () {
                        // TODO: Like action
                      },
                    ),
                    Text('${post['like_count'] ?? 0}'),
                    const SizedBox(width: 16),
                    IconButton(
                      icon: Icon(
                        Icons.mode_comment_outlined,
                        color: colorScheme.primary,
                      ),
                      onPressed: () {
                        context.go('/post/$postId');
                      },
                    ),
                    Text('${post['comment_count'] ?? 0}'),
                    const Spacer(),
                    IconButton(
                      icon: Icon(Icons.share, color: colorScheme.primary, size: 22),
                      onPressed: () async {
                        final url = 'https://www.thinkshare.com/post/$postId';
                        await Clipboard.setData(ClipboardData(text: url));
                        if (context.mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(
                              content: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: const [
                                  Icon(Icons.check_circle, color: Colors.green, size: 18),
                                  SizedBox(width: 8),
                                  Text('Url copied'),
                                ],
                              ),
                              backgroundColor: Colors.grey[900],
                              behavior: SnackBarBehavior.floating,
                              margin: const EdgeInsets.only(
                                bottom: 60, right: 20, left: 20,
                              ),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(8),
                              ),
                              duration: const Duration(seconds: 1),
                              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                            ),
                          );
                        }
                      },
                    ),
                  ],
                ),
              ),
            ],
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
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      padding: const EdgeInsets.all(16),
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
