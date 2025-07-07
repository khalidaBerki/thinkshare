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

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Comments',
          style: TextStyle(
            fontFamily: 'Montserrat',
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: Theme.of(context).colorScheme.primary,
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
                  color: Theme.of(
                    context,
                  ).colorScheme.onSurface.withOpacity(0.7),
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
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: controller,
                decoration: const InputDecoration(hintText: "Add a comment..."),
              ),
            ),
            IconButton(
              icon: const Icon(Icons.send),
              onPressed: () {
                if (controller.text.trim().isNotEmpty) {
                  onAddComment(controller.text.trim());
                  controller.clear();
                }
              },
            ),
          ],
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

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(backgroundImage: NetworkImage(avatarUrl), radius: 18),
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
                        onPressed: () =>
                            Navigator.pop(context, editController.text.trim()),
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
          ),
        ],
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
