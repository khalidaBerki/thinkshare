import 'package:flutter/material.dart';

class CommentSection extends StatelessWidget {
  final List<Map<String, dynamic>> comments;
  final bool isLoading;
  final Future<void> Function(String text) onAddComment;
  final Future<void> Function(String commentId, String text) onEditComment;
  final Future<void> Function(String commentId) onDeleteComment;

  const CommentSection({
    super.key,
    required this.comments,
    required this.isLoading,
    required this.onAddComment,
    required this.onEditComment,
    required this.onDeleteComment,
  });

  @override
  Widget build(BuildContext context) {
    final controller = TextEditingController();
    final colorScheme = Theme.of(context).colorScheme;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Comments',
          style: TextStyle(
            fontFamily: 'Montserrat',
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: colorScheme.primary,
          ),
        ),
        const SizedBox(height: 10),
        if (isLoading)
          const Center(child: CircularProgressIndicator())
        else if (comments.isEmpty)
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 20),
            child: Center(
              child: Text(
                'No comments yet.',
                style: TextStyle(
                  color: colorScheme.onSurface.withOpacity(0.7),
                  fontSize: 16,
                ),
              ),
            ),
          )
        else
          ...comments.map(
            (comment) => _CommentTile(
              comment: comment,
              onEdit: onEditComment,
              onDelete: onDeleteComment,
            ),
          ),
        const SizedBox(height: 10),
        Material(
          elevation: 3,
          borderRadius: BorderRadius.circular(24),
          color: colorScheme.surface,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: controller,
                    decoration: const InputDecoration(
                      hintText: "Add a comment...",
                      border: InputBorder.none,
                      contentPadding: EdgeInsets.symmetric(
                        horizontal: 10,
                        vertical: 12,
                      ),
                    ),
                  ),
                ),
                IconButton(
                  icon: Icon(Icons.send, color: colorScheme.primary),
                  splashRadius: 22,
                  onPressed: () {
                    if (controller.text.trim().isNotEmpty) {
                      onAddComment(controller.text.trim());
                      controller.clear();
                    }
                  },
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}

class _CommentTile extends StatelessWidget {
  final Map<String, dynamic> comment;
  final Future<void> Function(String commentId, String text) onEdit;
  final Future<void> Function(String commentId) onDelete;

  const _CommentTile({
    required this.comment,
    required this.onEdit,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final userName = comment['username'] ?? 'User';
    final avatarUrl =
        (comment['avatar_url'] != null &&
            comment['avatar_url'].toString().isNotEmpty)
        ? comment['avatar_url']
        : 'https://ui-avatars.com/api/?name=${Uri.encodeComponent(userName)}';

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Material(
        elevation: 2,
        borderRadius: BorderRadius.circular(16),
        color: colorScheme.surface,
        child: Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(16)),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              CircleAvatar(
                backgroundImage: NetworkImage(avatarUrl),
                radius: 18,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      userName,
                      style: TextStyle(
                        fontFamily: 'Montserrat',
                        fontWeight: FontWeight.bold,
                        color: colorScheme.primary,
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      comment['text'] ?? '',
                      style: TextStyle(
                        fontFamily: 'Montserrat',
                        fontSize: 14,
                        color: colorScheme.onSurface,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      _formatDate(comment['created_at']),
                      style: TextStyle(
                        color: colorScheme.secondary.withOpacity(0.7),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              PopupMenuButton<String>(
                onSelected: (value) async {
                  if (value == 'edit') {
                    final editController = TextEditingController(
                      text: comment['text'],
                    );
                    final result = await showDialog<String>(
                      context: context,
                      builder: (context) => AlertDialog(
                        title: const Text('Edit comment'),
                        content: TextField(
                          controller: editController,
                          autofocus: true,
                        ),
                        actions: [
                          TextButton(
                            onPressed: () => Navigator.pop(context),
                            child: const Text('Cancel'),
                          ),
                          TextButton(
                            onPressed: () => Navigator.pop(
                              context,
                              editController.text.trim(),
                            ),
                            child: const Text('Save'),
                          ),
                        ],
                      ),
                    );
                    if (result != null && result.isNotEmpty) {
                      await onEdit(comment['id'].toString(), result);
                    }
                  } else if (value == 'delete') {
                    await onDelete(comment['id'].toString());
                  }
                },
                itemBuilder: (context) => [
                  const PopupMenuItem(value: 'edit', child: Text('Edit')),
                  const PopupMenuItem(value: 'delete', child: Text('Delete')),
                ],
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
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
