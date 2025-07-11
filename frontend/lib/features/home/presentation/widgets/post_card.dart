import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/services.dart';
import 'media_carousel.dart';
import '../../../../config/api_config.dart';
import 'package:provider/provider.dart';
import '../../presentation/providers/home_provider.dart';
import '../../../../services/payment_service.dart';
import 'upgrade_banner.dart';

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
    final hasAccess = post['has_access'] == true;
    final isPaidOnly = post['is_paid_only'] == true;
    final monthlyPrice = creator['monthly_price'];
    final username = creator['username'];

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
      child: Material(
        elevation: 3,
        borderRadius: BorderRadius.circular(22),
        color: colorScheme.surface,
        child: InkWell(
          borderRadius: BorderRadius.circular(22),
          onTap: () {
            context.go('/post/$postId');
          },
          child: Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(22),
              boxShadow: [
                BoxShadow(
                  color: colorScheme.primary.withOpacity(0.06),
                  blurRadius: 12,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
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
                              ? '${ApiConfig.baseUrl}${creator['avatar_url']}'
                              : 'https://ui-avatars.com/api/?name=${Uri.encodeComponent(creator['full_name'] ?? 'User')}',
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
                  ),
                  if (!hasAccess && isPaidOnly && monthlyPrice != null && monthlyPrice > 0)
                    UpgradeBanner(
                      creatorId: creator['id'],
                      monthlyPrice: monthlyPrice is num ? monthlyPrice.toDouble() : null,
                      username: username,
                    )
                  else ...[
                    Padding(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 18.0,
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
                    if (mediaUrls.isNotEmpty)
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 8.0),
                        child: ClipRRect(
                          borderRadius: BorderRadius.circular(18),
                          child: MediaCarousel(mediaUrls: mediaUrls),
                        ),
                      ),
                  ],
                  Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12.0,
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
                          onPressed: () async {
                            final provider = Provider.of<HomeProvider>(
                              context,
                              listen: false,
                            );
                            await provider.toggleLike(post['id'].toString());
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
                          icon: Icon(
                            Icons.share,
                            color: colorScheme.primary,
                            size: 22,
                          ),
                          onPressed: () async {
                            final url =
                                'https://www.thinkshare.com/post/$postId';
                            await Clipboard.setData(ClipboardData(text: url));
                            if (context.mounted) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: const [
                                      Icon(
                                        Icons.check_circle,
                                        color: Colors.green,
                                        size: 18,
                                      ),
                                      SizedBox(width: 8),
                                      Text('Url copied'),
                                    ],
                                  ),
                                  backgroundColor: Colors.grey[900],
                                  behavior: SnackBarBehavior.floating,
                                  margin: const EdgeInsets.only(
                                    bottom: 60,
                                    right: 20,
                                    left: 20,
                                  ),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  duration: const Duration(seconds: 1),
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 16,
                                    vertical: 8,
                                  ),
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
