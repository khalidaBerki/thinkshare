import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../providers/message_provider.dart';

class ConversationScreen extends StatefulWidget {
  final int otherUserId;

  const ConversationScreen({super.key, required this.otherUserId});

  @override
  State<ConversationScreen> createState() => _ConversationScreenState();
}

class _ConversationScreenState extends State<ConversationScreen>
    with TickerProviderStateMixin {
  final TextEditingController _controller = TextEditingController();
  final GlobalKey<AnimatedListState> _listKey = GlobalKey<AnimatedListState>();
  final ScrollController _scrollController = ScrollController();
  int _lastMessageCount = 0;
  static const double _inputBarHeight = 64; // Height of the input bar (approx)

  @override
  void initState() {
    super.initState();
    Future.microtask(() async {
      await Provider.of<MessageProvider>(
        context,
        listen: false,
      ).fetchMessages(widget.otherUserId);
      setState(() {
        _lastMessageCount = Provider.of<MessageProvider>(
          context,
          listen: false,
        ).messages.length;
      });
    });
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final provider = Provider.of<MessageProvider>(context, listen: false);
    provider.addListener(_onMessagesChanged);
  }

  @override
  void dispose() {
    Provider.of<MessageProvider>(
      context,
      listen: false,
    ).removeListener(_onMessagesChanged);
    _controller.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _onMessagesChanged() {
    final provider = Provider.of<MessageProvider>(context, listen: false);
    if (provider.messages.length > _lastMessageCount) {
      _listKey.currentState?.insertItem(
        0,
        duration: const Duration(milliseconds: 400),
      );
      _lastMessageCount = provider.messages.length;

      // Scroll to top (latest message with reverse: true)
      Future.delayed(const Duration(milliseconds: 100), () {
        if (_scrollController.hasClients) {
          _scrollController.animateTo(
            0.0,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
          );
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<MessageProvider>();
    final messages = provider.messages;
    final isSending = provider.isSending;
    final colorScheme = Theme.of(context).colorScheme;

    // Get the other user from the first message
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
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () {
            context.go('/messages');
          },
        ),
        title: Row(
          children: [
            CircleAvatar(
              backgroundImage: avatarUrl.isNotEmpty
                  ? NetworkImage(avatarUrl)
                  : null,
              child: avatarUrl.isEmpty ? const Icon(Icons.person) : null,
            ),
            const SizedBox(width: 12),
            Text(username, style: const TextStyle(fontWeight: FontWeight.bold)),
          ],
        ),
        backgroundColor: colorScheme.primary,
        elevation: 1,
      ),
      body: Stack(
        children: [
          Container(
            color: colorScheme.surface,
            child: messages.isEmpty
                ? const Center(
                    child: Text(
                      "No messages yet.",
                      style: TextStyle(color: Colors.grey),
                    ),
                  )
                : AnimatedList(
                    key: _listKey,
                    controller: _scrollController,
                    reverse: true,
                    padding: EdgeInsets.fromLTRB(
                      8,
                      16,
                      8,
                      _inputBarHeight + 16, // Add bottom padding for input bar
                    ),
                    initialItemCount: messages.length,
                    physics: const BouncingScrollPhysics(),
                    itemBuilder: (context, i, animation) {
                      final msg = messages[messages.length - 1 - i];
                      final isMe = msg['sender']['id'].toString() != widget.otherUserId.toString();
                      final msgAvatar = isMe ? null : (otherUser['avatar_url'] ?? '');
                      final myAvatar = isMe ? (msg['sender']['avatar_url'] ?? '') : null;

                      return SlideTransition(
                        position: Tween<Offset>(
                          begin: Offset(isMe ? 1 : -1, 0),
                          end: Offset.zero,
                        ).animate(CurvedAnimation(
                          parent: animation,
                          curve: Curves.easeOutBack,
                        )),
                        child: FadeTransition(
                          opacity: animation,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 4), // Plus d'espace entre bulles
                            child: Row(
                              mainAxisAlignment: isMe ? MainAxisAlignment.end : MainAxisAlignment.start,
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                if (!isMe)
                                  Padding(
                                    padding: const EdgeInsets.only(right: 6.0),
                                    child: CircleAvatar(
                                      radius: 18,
                                      backgroundImage: msgAvatar != null && msgAvatar.isNotEmpty
                                          ? NetworkImage(msgAvatar)
                                          : null,
                                      child: (msgAvatar == null || msgAvatar.isEmpty)
                                          ? const Icon(Icons.person, size: 18)
                                          : null,
                                    ),
                                  ),
                                // Limite la largeur max de la bulle
                                ConstrainedBox(
                                  constraints: BoxConstraints(
                                    maxWidth: MediaQuery.of(context).size.width * 0.75,
                                  ),
                                  child: AnimatedContainer(
                                    duration: const Duration(milliseconds: 250),
                                    curve: Curves.easeInOut,
                                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                                    decoration: BoxDecoration(
                                      gradient: isMe
                                          ? LinearGradient(
                                              colors: [
                                                colorScheme.primary,
                                                colorScheme.primary.withOpacity(0.8),
                                              ],
                                              begin: Alignment.topRight,
                                              end: Alignment.bottomLeft,
                                            )
                                          : null,
                                      color: isMe ? null : colorScheme.surfaceContainerHighest,
                                      borderRadius: BorderRadius.only(
                                        topLeft: const Radius.circular(18),
                                        topRight: const Radius.circular(18),
                                        bottomLeft: Radius.circular(isMe ? 18 : 4),
                                        bottomRight: Radius.circular(isMe ? 4 : 18),
                                      ),
                                      boxShadow: [
                                        BoxShadow(
                                          color: Colors.black.withOpacity(0.06),
                                          blurRadius: 8,
                                          offset: const Offset(0, 2),
                                        ),
                                      ],
                                    ),
                                    child: Column(
                                      crossAxisAlignment: isMe
                                          ? CrossAxisAlignment.end
                                          : CrossAxisAlignment.start,
                                      children: [
                                        Row(
                                          mainAxisSize: MainAxisSize.min,
                                          children: [
                                            Flexible(
                                              child: Text(
                                                msg['content'] ?? '',
                                                style: TextStyle(
                                                  color: isMe
                                                      ? Colors.white
                                                      : colorScheme.onSurface,
                                                  fontSize: 16,
                                                  fontWeight: FontWeight.w500,
                                                ),
                                              ),
                                            ),
                                            if (isMe)
                                              PopupMenuButton<String>(
                                                icon: Icon(Icons.more_vert,
                                                    size: 18,
                                                    color: isMe
                                                        ? Colors.white70
                                                        : Colors.grey),
                                                shape: RoundedRectangleBorder(
                                                    borderRadius: BorderRadius.circular(12)),
                                                onSelected: (value) async {
                                                  if (value == 'edit') {
                                                    final controller = TextEditingController(
                                                        text: msg['content']);
                                                    final result = await showDialog<String>(
                                                      context: context,
                                                      builder: (context) => AlertDialog(
                                                        title: const Text('Edit message'),
                                                        content: TextField(
                                                          controller: controller,
                                                          autofocus: true,
                                                          decoration: const InputDecoration(
                                                            border: OutlineInputBorder(),
                                                          ),
                                                          minLines: 1,
                                                          maxLines: 4,
                                                        ),
                                                        shape: RoundedRectangleBorder(
                                                            borderRadius:
                                                                BorderRadius.circular(18)),
                                                        actions: [
                                                          TextButton(
                                                            onPressed: () =>
                                                                Navigator.pop(context),
                                                            child: const Text('Cancel'),
                                                          ),
                                                          ElevatedButton(
                                                            onPressed: () => Navigator.pop(
                                                                context, controller.text.trim()),
                                                            style: ElevatedButton.styleFrom(
                                                              backgroundColor: Theme.of(context)
                                                                  .colorScheme
                                                                  .primary,
                                                              foregroundColor: Colors.white,
                                                              shape: RoundedRectangleBorder(
                                                                  borderRadius:
                                                                      BorderRadius.circular(8)),
                                                            ),
                                                            child: const Text('Save'),
                                                          ),
                                                        ],
                                                      ),
                                                    );
                                                    if (result != null &&
                                                        result.isNotEmpty &&
                                                        result != msg['content']) {
                                                      await Provider.of<MessageProvider>(context,
                                                              listen: false)
                                                          .updateMessage(msg['id'], result);
                                                    }
                                                  } else if (value == 'delete') {
                                                    final confirm = await showDialog<bool>(
                                                      context: context,
                                                      builder: (context) => AlertDialog(
                                                        title: const Text('Delete message'),
                                                        content: const Text(
                                                            'Are you sure you want to delete this message?'),
                                                        shape: RoundedRectangleBorder(
                                                            borderRadius:
                                                                BorderRadius.circular(18)),
                                                        actions: [
                                                          TextButton(
                                                            onPressed: () =>
                                                                Navigator.pop(context, false),
                                                            child: const Text('Cancel'),
                                                          ),
                                                          ElevatedButton(
                                                            onPressed: () =>
                                                                Navigator.pop(context, true),
                                                            style: ElevatedButton.styleFrom(
                                                              backgroundColor: Theme.of(context)
                                                                  .colorScheme
                                                                  .error,
                                                              foregroundColor: Colors.white,
                                                              shape: RoundedRectangleBorder(
                                                                  borderRadius:
                                                                      BorderRadius.circular(8)),
                                                            ),
                                                            child: const Text('Delete'),
                                                          ),
                                                        ],
                                                      ),
                                                    );
                                                    if (confirm == true) {
                                                      await Provider.of<MessageProvider>(context,
                                                              listen: false)
                                                          .deleteMessage(msg['id']);
                                                    }
                                                  }
                                                },
                                                itemBuilder: (context) => [
                                                  const PopupMenuItem(
                                                      value: 'edit', child: Text('Edit')),
                                                  const PopupMenuItem(
                                                      value: 'delete', child: Text('Delete')),
                                                ],
                                              ),
                                          ],
                                        ),
                                        const SizedBox(height: 4),
                                        Row(
                                          mainAxisSize: MainAxisSize.min,
                                          children: [
                                            Text(
                                              msg['created_at']?.substring(11, 16) ?? '',
                                              style: TextStyle(
                                                fontSize: 11,
                                                color: isMe
                                                    ? Colors.white70
                                                    : Colors.grey,
                                              ),
                                            ),
                                            if (isMe && msg['status'] == 'READ')
                                              const Padding(
                                                padding: EdgeInsets.only(left: 4),
                                                child: Icon(
                                                  Icons.done_all,
                                                  size: 16,
                                                  color: Colors.lightBlueAccent,
                                                ),
                                              ),
                                          ],
                                        ),
                                      ],
                                    ),
                                  ),
                                ),
                                // Avatar right (me)
                                if (isMe)
                                  Padding(
                                    padding: const EdgeInsets.only(left: 6.0),
                                    child: CircleAvatar(
                                      radius: 18,
                                      backgroundImage: myAvatar != null && myAvatar.isNotEmpty
                                          ? NetworkImage(myAvatar)
                                          : null,
                                      child: (myAvatar == null || myAvatar.isEmpty)
                                          ? const Icon(Icons.person, size: 18)
                                          : null,
                                    ),
                                  ),
                              ],
                            ),
                          ),
                        ),
                      );
                    },
                  ),
          ),
          // Floating input bar
          Positioned(
            left: 0,
            right: 0,
            bottom: 0,
            child: SafeArea(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                child: Material(
                  elevation: 8,
                  borderRadius: BorderRadius.circular(32),
                  color: colorScheme.surface,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 6,
                    ),
                    child: Row(
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _controller,
                            decoration: InputDecoration(
                              hintText: "Type your message...",
                              border: InputBorder.none,
                              contentPadding: const EdgeInsets.symmetric(
                                horizontal: 8,
                                vertical: 10,
                              ),
                            ),
                            minLines: 1,
                            maxLines: 4,
                          ),
                        ),
                        AnimatedSwitcher(
                          duration: const Duration(milliseconds: 200),
                          child: isSending
                              ? const Padding(
                                  padding: EdgeInsets.all(8.0),
                                  child: SizedBox(
                                    width: 24,
                                    height: 24,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                    ),
                                  ),
                                )
                              : IconButton(
                                  icon: const Icon(Icons.send, size: 28),
                                  color: colorScheme.primary,
                                  onPressed: () async {
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
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
