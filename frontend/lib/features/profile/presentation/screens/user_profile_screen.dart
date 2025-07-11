import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';

import '../providers/profile_provider.dart';

class UserProfileScreen extends StatefulWidget {
  final int userId;
  const UserProfileScreen({super.key, required this.userId});

  @override
  State<UserProfileScreen> createState() => _UserProfileScreenState();
}

class _UserProfileScreenState extends State<UserProfileScreen> {
  @override
  void initState() {
    super.initState();
    final provider = Provider.of<ProfileProvider>(context, listen: false);
    provider.fetchUserProfile(widget.userId);
    provider.fetchUserPosts(widget.userId);
  }

  @override
  Widget build(BuildContext context) {
    final provider = Provider.of<ProfileProvider>(context);
    final colorScheme = Theme.of(context).colorScheme;

    if (provider.isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    if (provider.error != null) {
      return Scaffold(
        body: Center(
          child: Text(
            provider.error!,
            style: const TextStyle(color: Colors.red),
          ),
        ),
      );
    }
    final profile = provider.userProfile;
    if (profile == null) {
      return const Scaffold(body: Center(child: Text("User not found")));
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text(
          "User Profil",
          style: TextStyle(
            fontFamily: 'Montserrat',
            fontWeight: FontWeight.bold,
          ),
        ),
        centerTitle: true,
        elevation: 2,
        backgroundColor: colorScheme.surface,
        shadowColor: colorScheme.primary.withOpacity(0.08),
        surfaceTintColor: colorScheme.primary,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(18),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            CircleAvatar(
              radius: 44,
              backgroundImage:
                  profile['avatar_url'] != null &&
                      profile['avatar_url'].toString().isNotEmpty
                  ? NetworkImage(profile['avatar_url'])
                  : null,
              child:
                  (profile['avatar_url'] == null ||
                      profile['avatar_url'].toString().isEmpty)
                  ? const Icon(Icons.person, size: 44)
                  : null,
            ),
            const SizedBox(height: 12),
            Text(
              profile['full_name'] ?? '',
              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 20),
            ),
            if (profile['bio'] != null && profile['bio'].toString().isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(top: 4),
                child: Text(
                  'bio : ${profile['bio']}',
                  style: const TextStyle(fontSize: 15, color: Colors.grey),
                ),
              ),
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                ElevatedButton(
                  onPressed: () {}, // TODO: implement follow
                  style: ElevatedButton.styleFrom(
                    backgroundColor: colorScheme.surface,
                    foregroundColor: colorScheme.primary,
                    side: BorderSide(color: colorScheme.primary),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    elevation: 0,
                  ),
                  child: const Text("Follow"),
                ),
                const SizedBox(width: 12),
                ElevatedButton(
                  onPressed: () {
                    context.go('/messages/${widget.userId}');
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: colorScheme.surface,
                    foregroundColor: colorScheme.primary,
                    side: BorderSide(color: colorScheme.primary),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    elevation: 0,
                  ),
                  child: const Text("Send message"),
                ),
              ],
            ),
            const SizedBox(height: 18),
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(vertical: 8),
              decoration: BoxDecoration(
                color: colorScheme.surfaceVariant.withOpacity(0.7),
                borderRadius: BorderRadius.circular(12),
              ),
              child: const Center(
                child: Text(
                  "User Posts",
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
              ),
            ),
            const SizedBox(height: 10),
            if (provider.isPostsLoading)
              const Center(child: CircularProgressIndicator())
            else if (provider.postsError != null)
              Center(
                child: Text(
                  provider.postsError!,
                  style: const TextStyle(color: Colors.red),
                ),
              )
            else if (provider.userPosts.isEmpty)
              const Center(child: Text("No posts yet"))
            else
              GridView.builder(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: provider.userPosts.length,
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 3,
                  crossAxisSpacing: 8,
                  mainAxisSpacing: 8,
                  childAspectRatio: 0.7,
                ),
                itemBuilder: (context, i) {
                  final post = provider.userPosts[i];
                  final hasMedia =
                      post['media_urls'] != null &&
                      post['media_urls'].isNotEmpty;
                  final mediaUrl = hasMedia
                      ? post['media_urls'][0].toString().replaceAll('\\', '/')
                      : null;
                  final isVideo =
                      mediaUrl != null &&
                      (mediaUrl.endsWith('.mp4') || mediaUrl.endsWith('.mov'));

                  return GestureDetector(
                    onTap: () {
                      context.go('/post/${post['id']}');
                    },
                    child: Container(
                      decoration: BoxDecoration(
                        color: colorScheme.surfaceVariant,
                        borderRadius: BorderRadius.circular(10),
                        border: Border.all(
                          color: colorScheme.primary, // Couleur de la bordure
                          width: 1.5,
                        ),
                      ),
                      child: hasMedia
                          ? ClipRRect(
                              borderRadius: BorderRadius.circular(10),
                              child: isVideo
                                  ? Stack(
                                      fit: StackFit.expand,
                                      children: [
                                        // Optionnel: ajouter une miniature vidéo ici
                                        Container(color: Colors.black12),
                                        const Center(
                                          child: Icon(
                                            Icons.play_circle_fill,
                                            size: 40,
                                            color: Colors.white70,
                                          ),
                                        ),
                                      ],
                                    )
                                  : Image.network(
                                      mediaUrl!,
                                      fit: BoxFit.cover,
                                      errorBuilder:
                                          (context, error, stackTrace) =>
                                              const Center(
                                                child: Icon(
                                                  Icons.broken_image,
                                                  size: 40,
                                                  color: Colors.grey,
                                                ),
                                              ),
                                    ),
                            )
                          : const Center(
                              child: Icon(
                                Icons.image,
                                size: 40,
                                color: Colors.grey,
                              ),
                            ),
                    ),
                  );
                },
              ),
          ],
        ),
      ),
    );
  }
}
