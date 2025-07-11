import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/home_provider.dart';
import '../widgets/post_card.dart';

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
              return PostCard(post: post);
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
}
