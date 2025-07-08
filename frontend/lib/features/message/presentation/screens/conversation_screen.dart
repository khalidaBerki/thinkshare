import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/message_provider.dart';

class ConversationScreen extends StatefulWidget {
  final int otherUserId;

  const ConversationScreen({super.key, required this.otherUserId});

  @override
  State<ConversationScreen> createState() => _ConversationScreenState();
}

class _ConversationScreenState extends State<ConversationScreen> {
  final TextEditingController _controller = TextEditingController();

  @override
  void initState() {
    super.initState();
    Future.microtask(
      () => Provider.of<MessageProvider>(
        context,
        listen: false,
      ).fetchMessages(widget.otherUserId),
    );
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<MessageProvider>();
    final messages = provider.messages;
    final isSending = provider.isSending;
    final colorScheme = Theme.of(context).colorScheme;

    // Après avoir fetch les messages :
    final firstMsg = messages.isNotEmpty ? messages.first : null;
    final otherUser = firstMsg != null
        ? (firstMsg['sender']['id'] == widget.otherUserId
              ? firstMsg['sender']
              : firstMsg['receiver'])
        : {};

    final username = otherUser['username'] ?? '';
    final avatarUrl = otherUser['avatar_url'] ?? '';

    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            CircleAvatar(backgroundImage: NetworkImage(avatarUrl)),
            const SizedBox(width: 12),
            Text(username),
          ],
        ),
        backgroundColor: colorScheme.primary,
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              reverse: true,
              padding: const EdgeInsets.all(12),
              itemCount: messages.length,
              itemBuilder: (context, i) {
                final msg = messages[messages.length - 1 - i];
                final isMe =
                    msg['sender']['id'].toString() !=
                    widget.otherUserId.toString();
                return Align(
                  alignment: isMe
                      ? Alignment.centerRight
                      : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.symmetric(vertical: 4),
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 8,
                    ),
                    decoration: BoxDecoration(
                      color: isMe
                          ? colorScheme.primary
                          : colorScheme.surfaceVariant,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      crossAxisAlignment: isMe
                          ? CrossAxisAlignment.end
                          : CrossAxisAlignment.start,
                      children: [
                        Text(
                          msg['content'] ?? '',
                          style: TextStyle(
                            color: isMe ? Colors.white : colorScheme.onSurface,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          msg['created_at']?.substring(11, 16) ?? '',
                          style: const TextStyle(
                            fontSize: 10,
                            color: Colors.grey,
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
          Divider(height: 1),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _controller,
                    decoration: const InputDecoration(
                      hintText: "Type message here...",
                      border: OutlineInputBorder(),
                      isDense: true,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: isSending
                      ? const CircularProgressIndicator()
                      : const Icon(Icons.send),
                  onPressed: isSending
                      ? null
                      : () async {
                          final text = _controller.text.trim();
                          if (text.isNotEmpty) {
                            await provider.sendMessage(
                              widget.otherUserId,
                              text,
                            );
                            _controller.clear();
                          }
                        },
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
