import 'package:flutter/material.dart';
import '../../data/message_repository.dart';
import '../../data/message_api.dart';
import 'package:go_router/go_router.dart';

class MessageProvider extends ChangeNotifier {
  final MessageRepository _repo = MessageRepository(MessageApi());

  List<Map<String, dynamic>> conversations = [];
  List<Map<String, dynamic>> messages = [];
  bool isLoading = false;
  bool isSending = false;

  MessageProvider();

  Future<void> fetchConversations() async {
    isLoading = true;
    notifyListeners();
    try {
      conversations = await _repo.getConversations();
    } catch (e) {
      conversations = [];
    }
    isLoading = false;
    notifyListeners();
  }

  Future<void> fetchMessages(int otherUserId) async {
    isLoading = true;
    notifyListeners();
    try {
      messages = await _repo.getMessagesWithUser(otherUserId);
      // Marquer comme lu après avoir récupéré les messages
      await markAsRead(otherUserId);
    } catch (e) {
      messages = [];
    }
    isLoading = false;
    notifyListeners();
  }

  Future<void> sendMessage(int receiverId, String content) async {
    isSending = true;
    notifyListeners();
    try {
      final msg = await _repo.sendMessage(receiverId, content);
      messages.add(msg);
    } finally {
      isSending = false;
      notifyListeners();
    }
  }

  Future<void> markAsRead(int senderId) async {
    await _repo.markAsRead(senderId);
  }

  Future<void> updateMessage(int messageId, String content) async {
    final updated = await _repo.updateMessage(messageId, content);
    final idx = messages.indexWhere((m) => m['id'] == messageId);
    if (idx != -1) {
      messages[idx] = updated;
      notifyListeners();
    }
  }

  Future<void> deleteMessage(int messageId) async {
    await _repo.deleteMessage(messageId);
    messages.removeWhere((m) => m['id'] == messageId);
    notifyListeners();
  }

  void openConversation(BuildContext context, Map<String, dynamic> conv) {
    context.go('/messages/${conv['other_user']['id']}');
  }
}
