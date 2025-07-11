import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../features/home/presentation/screens/feed_screen.dart';
import '../../features/home/presentation/screens/create_post_screen.dart';
import '../../features/message/presentation/screens/conversation_list_screen.dart';
import '../../features/profile/presentation/screens/my_profile_screen.dart';

class NavigationScreen extends StatefulWidget {
  final int tabIndex;
  const NavigationScreen({super.key, this.tabIndex = 0});

  @override
  State<NavigationScreen> createState() => _NavigationScreenState();
}

class _NavigationScreenState extends State<NavigationScreen> {
  late int _currentIndex;
  late final PageController _pageController;

  @override
  void initState() {
    super.initState();
    _currentIndex = widget.tabIndex;
    _pageController = PageController(initialPage: _currentIndex);
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void onPageChanged(int page) {
    setState(() {
      _currentIndex = page;
    });
    switch (page) {
      case 0:
        context.go('/home');
        break;
      case 1:
        context.go('/post');
        break;
      case 2:
        context.go('/messages');
        break;
      case 3:
        context.go('/profile');
        break;
    }
  }

  void onNavTapped(int page) {
    _pageController.jumpToPage(page);
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Scaffold(
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        selectedItemColor: colorScheme.primary,
        unselectedItemColor: colorScheme.secondary,
        backgroundColor: colorScheme.surface,
        currentIndex: _currentIndex,
        onTap: onNavTapped,
        showSelectedLabels: false,
        showUnselectedLabels: false,
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home_outlined), label: ''),
          BottomNavigationBarItem(
            icon: Icon(Icons.add_box_outlined),
            label: '',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.chat_bubble_outline),
            label: '',
          ),
          BottomNavigationBarItem(icon: Icon(Icons.person_outline), label: ''),
        ],
      ),
      body: PageView(
        controller: _pageController,
        onPageChanged: onPageChanged,
        physics: const BouncingScrollPhysics(),
        children: const [
          FeedScreen(),
          CreatePostScreen(),
          ConversationListScreen(),
          MyProfileScreen(),
        ],
      ),
    );
  }
}
