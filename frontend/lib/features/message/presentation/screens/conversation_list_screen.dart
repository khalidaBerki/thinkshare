import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/message_provider.dart';
import '../widgets/conversation_tile.dart';

class ConversationListScreen extends StatefulWidget {
  const ConversationListScreen({super.key});

  @override
  State<ConversationListScreen> createState() => _ConversationListScreenState();
}

class _ConversationListScreenState extends State<ConversationListScreen> {
  @override
  void initState() {
    super.initState();
    Future.microtask(
      () => Provider.of<MessageProvider>(
        context,
        listen: false,
      ).fetchConversations(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<MessageProvider>();
    final conversations = provider.conversations;
    final isLoading = provider.isLoading;
    final isWide = MediaQuery.of(context).size.width > 600;
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text("Messages"),
        centerTitle: true,
        backgroundColor: colorScheme.primary,
      ),
      body: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: isWide ? 500 : double.infinity),
          child: isLoading
              ? const Center(child: CircularProgressIndicator())
              : RefreshIndicator(
                  onRefresh: provider.fetchConversations,
                  child: ListView.separated(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8,
                    ),
                    itemCount: conversations.length,
                    separatorBuilder: (_, __) => const Divider(height: 1),
                    itemBuilder: (context, i) => ConversationTile(
                      conversation: conversations[i],
                      onTap: () =>
                          provider.openConversation(context, conversations[i]),
                    ),
                  ),
                ),
        ),
      ),
    );
  }
}
