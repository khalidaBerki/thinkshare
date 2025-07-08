import 'package:flutter/material.dart';

class ConversationTile extends StatelessWidget {
  final Map<String, dynamic> conversation;
  final VoidCallback onTap;

  const ConversationTile({
    super.key,
    required this.conversation,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final other = conversation['other_user'] ?? {};
    final unread = conversation['unread_count'] ?? 0;
    final lastMsg = conversation['last_message'] ?? '';
    final timestamp = conversation['timestamp'] ?? '';
    final colorScheme = Theme.of(context).colorScheme;

    return Material(
      color: unread > 0
          ? colorScheme.primary.withOpacity(0.07)
          : Colors.transparent,
      borderRadius: BorderRadius.circular(12),
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: onTap,
        hoverColor: colorScheme.primary.withOpacity(0.12),
        child: ListTile(
          leading: CircleAvatar(
            backgroundImage: (other['avatar_url'] ?? '').isNotEmpty
                ? NetworkImage(other['avatar_url'])
                : const AssetImage('assets/images/icon.png') as ImageProvider,
            radius: 26,
          ),
          title: Text(
            other['username'] ?? '',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 16,
              color: colorScheme.onSurface,
            ),
          ),
          subtitle: Padding(
            padding: const EdgeInsets.only(top: 2.0),
            child: Text(
              lastMsg,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                color: unread > 0
                    ? colorScheme.primary
                    : colorScheme.onSurface.withOpacity(0.7),
                fontWeight: unread > 0 ? FontWeight.bold : FontWeight.normal,
                fontSize: 14,
              ),
            ),
          ),
          trailing: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                _formatTime(timestamp),
                style: const TextStyle(fontSize: 12, color: Colors.grey),
              ),
              if (unread > 0)
                Container(
                  margin: const EdgeInsets.only(top: 4),
                  padding: const EdgeInsets.symmetric(
                    horizontal: 10,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: colorScheme.primary,
                    borderRadius: BorderRadius.circular(20),
                    boxShadow: [
                      BoxShadow(
                        color: colorScheme.primary.withOpacity(0.2),
                        blurRadius: 4,
                        offset: const Offset(0, 2),
                      ),
                    ],
                  ),
                  child: Text(
                    unread.toString(),
                    semanticsLabel: 'Unread messages',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 13,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
            ],
          ),
          contentPadding: const EdgeInsets.symmetric(
            horizontal: 16,
            vertical: 6,
          ),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
    );
  }

  String _formatTime(String iso) {
    if (iso.isEmpty) return '';
    final dt = DateTime.tryParse(iso);
    if (dt == null) return '';
    final now = DateTime.now();
    if (now.difference(dt).inDays == 0) {
      // Aujourd'hui
      return "${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}";
    } else if (now.difference(dt).inDays == 1) {
      return "Yesterday";
    } else {
      return "${dt.day}/${dt.month}/${dt.year}";
    }
  }
}
