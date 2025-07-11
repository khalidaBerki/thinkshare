import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';
import 'dart:io';
import '../providers/profile_provider.dart';

class MyProfileScreen extends StatefulWidget {
  const MyProfileScreen({super.key});

  @override
  State<MyProfileScreen> createState() => _MyProfileScreenState();
}

class _MyProfileScreenState extends State<MyProfileScreen> {
  final _formKey = GlobalKey<FormState>();
  String? _fullName;
  String? _bio;
  String? _avatarUrl;
  File? _avatarFile;

  Future<void> _showEditDialog(
    BuildContext context,
    Map<String, dynamic> profile,
  ) async {
    _fullName = profile['full_name'];
    _bio = profile['bio'];
    _avatarUrl = profile['avatar_url'];
    _avatarFile = null;

    await showDialog(
      context: context,
      builder: (context) {
        final provider = Provider.of<ProfileProvider>(context);
        return AlertDialog(
          title: const Text("Edit my profile"),
          content: Form(
            key: _formKey,
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  GestureDetector(
                    onTap: () async {
                      final picker = ImagePicker();
                      final picked = await picker.pickImage(
                        source: ImageSource.gallery,
                      );
                      if (picked != null) {
                        setState(() {
                          _avatarFile = File(picked.path);
                          _avatarUrl = null;
                        });
                      }
                    },
                    child: CircleAvatar(
                      radius: 36,
                      backgroundImage: _avatarFile != null
                          ? FileImage(_avatarFile!)
                          : (_avatarUrl != null && _avatarUrl!.isNotEmpty
                                    ? NetworkImage(_avatarUrl!)
                                    : null)
                                as ImageProvider?,
                      child:
                          (_avatarFile == null &&
                              (_avatarUrl == null || _avatarUrl!.isEmpty))
                          ? const Icon(Icons.camera_alt, size: 36)
                          : null,
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    initialValue: _fullName,
                    decoration: const InputDecoration(labelText: "Full name"),
                    onChanged: (v) => _fullName = v,
                  ),
                  const SizedBox(height: 8),
                  TextFormField(
                    initialValue: _bio,
                    decoration: const InputDecoration(labelText: "Bio"),
                    maxLines: 2,
                    onChanged: (v) => _bio = v,
                  ),
                ],
              ),
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text("Cancel"),
            ),
            ElevatedButton(
              onPressed: provider.isUpdating
                  ? null
                  : () async {
                      await provider.updateProfile(
                        fullName: _fullName,
                        bio: _bio,
                        avatarUrl: _avatarUrl,
                      );
                      if (provider.updateError == null) {
                        Navigator.pop(context);
                      }
                    },
              child: provider.isUpdating
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Text("Save"),
            ),
          ],
        );
      },
    );
  }

  @override
  void initState() {
    super.initState();
    final provider = Provider.of<ProfileProvider>(context, listen: false);
    provider.fetchMyProfile().then((_) {
      if (provider.myProfile != null) {
        provider.fetchMyPosts();
      }
    });
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
    final profile = provider.myProfile;
    if (profile == null) {
      return const Scaffold(body: Center(child: Text("Profile not found")));
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text(
          "My Profile",
          style: TextStyle(
            fontFamily: 'Montserrat',
            fontWeight: FontWeight.bold,
          ),
        ),
        centerTitle: true,
        elevation: 2,
        backgroundColor: colorScheme.surface,
        actions: [
          IconButton(
            icon: const Icon(Icons.dashboard),
            onPressed: () => context.go('/dashboard'),
            tooltip: "Dashboard",
          ),
          IconButton(
            icon: const Icon(Icons.logout, color: Colors.red),
            onPressed: () => provider.logout(context),
            tooltip: "Logout",
          ),
        ],
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
                  'bio: ${profile['bio']}',
                  style: const TextStyle(fontSize: 15, color: Colors.grey),
                ),
              ),
            const SizedBox(height: 10),
            ElevatedButton(
              onPressed: () => _showEditDialog(context, profile),
              style: ElevatedButton.styleFrom(
                backgroundColor: colorScheme.surface,
                foregroundColor: colorScheme.primary,
                side: BorderSide(color: colorScheme.primary),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
                elevation: 0,
              ),
              child: const Text("Edit my profile"),
            ),
            if (provider.updateError != null)
              Padding(
                padding: const EdgeInsets.only(top: 8),
                child: Text(
                  provider.updateError!,
                  style: const TextStyle(color: Colors.red),
                ),
              ),
            if (provider.updateSuccess != null)
              Padding(
                padding: const EdgeInsets.only(top: 8),
                child: Text(
                  provider.updateSuccess!,
                  style: TextStyle(color: colorScheme.primary),
                ),
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
                  "My Posts",
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
            else if (provider.myPosts.isEmpty)
              const Center(child: Text("No posts yet"))
            else
              GridView.builder(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: provider.myPosts.length,
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 3,
                  crossAxisSpacing: 8,
                  mainAxisSpacing: 8,
                  childAspectRatio: 0.7,
                ),
                itemBuilder: (context, i) {
                  final post = provider.myPosts[i];
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
