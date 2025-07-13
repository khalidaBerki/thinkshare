import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../providers/home_provider.dart';
import '../widgets/post_card.dart';
import '../widgets/upgrade_banner.dart';

class FeedScreen extends StatefulWidget {
  const FeedScreen({super.key});

  @override
  State<FeedScreen> createState() => _FeedScreenState();
}

class _FeedScreenState extends State<FeedScreen> {
  late ScrollController _controller;

  @override
  void initState() {
    super.initState();
    final provider = Provider.of<HomeProvider>(context, listen: false);
    provider.loadPosts();
    _controller = ScrollController()..addListener(_onScroll);
  }

  void _onScroll() {
    final provider = Provider.of<HomeProvider>(context, listen: false);
    if (_controller.position.pixels >=
        _controller.position.maxScrollExtent - 300) {
      if (provider.hasMore && !provider.isLoading) {
        provider.loadPosts();
      }
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final provider = Provider.of<HomeProvider>(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Home Feed'), centerTitle: true),
      backgroundColor: Theme.of(context).colorScheme.surface,
      body: RefreshIndicator(
        onRefresh: () => provider.loadPosts(refresh: true),
        child: ListView.builder(
          controller: _controller,
          physics: const BouncingScrollPhysics(),
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 14),
          itemCount: provider.posts.length + (provider.hasMore ? 1 : 0),
          itemBuilder: (context, index) {
            if (index < provider.posts.length) {
              final post = provider.posts[index];
              final hasAccess = post['has_access'] == true;
              final monthlyPrice = post['creator']?['monthly_price'];
              final isPaidOnly = post['is_paid_only'] == true;
              if (!hasAccess && isPaidOnly && monthlyPrice != null && monthlyPrice > 0) {
                final creator = post['creator'] ?? {};
                return Padding(
                  padding: const EdgeInsets.symmetric(vertical: 8.0),
                  child: Card(
                    elevation: 2,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(18),
                    ),
                    child: InkWell(
                      borderRadius: BorderRadius.circular(18),
                      onTap: () {
                        final postId = post['id'].toString();
                        context.go('/post/$postId');
                      },
                      child: Padding(
                        padding: const EdgeInsets.all(10.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            ListTile(
                              leading: CircleAvatar(
                                backgroundImage: NetworkImage(
                                  creator['avatar_url']?.isNotEmpty == true
                                      ? 'https://www.thinkshare.com/${creator['avatar_url']}'
                                      : 'https://ui-avatars.com/api/?name=${creator['full_name'] ?? 'User'}',
                                ),
                                radius: 22,
                              ),
                              title: Text(
                                creator['full_name'] ?? 'No Name',
                                style: const TextStyle(fontWeight: FontWeight.bold),
                              ),
                              subtitle: Text(_formatDate(post['created_at'])),
                            ),
                            Row(
                              children: [
                                Icon(Icons.star, color: Theme.of(context).colorScheme.primary, size: 20),
                                const SizedBox(width: 4),
                                Text('${post['like_count'] ?? 0}'),
                                const SizedBox(width: 16),
                                Icon(Icons.mode_comment_outlined, color: Theme.of(context).colorScheme.primary, size: 20),
                                const SizedBox(width: 4),
                                Text('${post['comment_count'] ?? 0}'),
                              ],
                            ),
                            const SizedBox(height: 8),
                            UpgradeBanner(
                              creatorId: (creator['id'] is int) ? creator['id'] : null,
                              monthlyPrice: monthlyPrice,
                              username: creator['username']?.toString() ?? '',
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                );
              } else {
                return PostCard(post: post);
              }
            } else {
              return const Center(
                child: Padding(
                  padding: EdgeInsets.all(24.0),
                  child: CircularProgressIndicator(),
                ),
              );
            }
          },
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
