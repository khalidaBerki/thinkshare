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
        title: const Text(
          "Messages",
          style: TextStyle(
            fontFamily: 'Montserrat',
            fontWeight: FontWeight.bold,
            fontSize: 20,
            letterSpacing: 0.2,
          ),
        ),
        centerTitle: true,
        elevation: 2,
        backgroundColor: colorScheme.surface,
        shadowColor: colorScheme.primary.withOpacity(0.08),
        surfaceTintColor: colorScheme.primary,
        automaticallyImplyLeading: false, // Pas de flèche retour
      ),
      body: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: isWide ? 500 : double.infinity),
          child: isLoading
              ? const Center(child: CircularProgressIndicator())
              : RefreshIndicator(
                  onRefresh: provider.fetchConversations,
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 12),
                    physics: const BouncingScrollPhysics(),
                    itemCount: conversations.length,
                    itemBuilder: (context, i) => Padding(
                      padding: const EdgeInsets.symmetric(vertical: 6),
                      child: ConversationTile(
                        conversation: conversations[i],
                        onTap: () =>
                            provider.openConversation(context, conversations[i]),
                      ),
                    ),
                  ),
                ),
        ),
      ),
    );
  }
}
