import 'package:go_router/go_router.dart';

import '../features/auth/presentation/screens/entry_screen.dart';
import '../features/auth/presentation/screens/login_screen.dart';
import '../features/auth/presentation/screens/register_screen.dart';
import '../features/home/presentation/screens/post_detail_screen.dart';
import '../features/message/presentation/screens/conversation_list_screen.dart';
import '../features/message/presentation/screens/conversation_screen.dart';
import '../core/widgets/navigation.dart';

final GoRouter appRouter = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(path: '/', builder: (context, state) => const EntryScreen()),

    GoRoute(path: '/login', builder: (context, state) => const LoginScreen()),

    GoRoute(
      path: '/register',
      builder: (context, state) => const RegisterScreen(),
    ),

    GoRoute(
      path: '/home',
      builder: (context, state) => const NavigationScreen(),
    ),

    GoRoute(
      path: '/messages',
      builder: (context, state) => const ConversationListScreen(),
    ),

    GoRoute(
      path: '/post/:id',
      builder: (context, state) {
        final id = state.pathParameters['id']!;
        return PostDetailScreen(postId: id);
      },
    ),
    GoRoute(
      path: '/messages/:otherUserId',
      builder: (context, state) {
        final otherUserId = int.parse(state.pathParameters['otherUserId']!);
        // Tu peux aussi passer username/avatarUrl via extra ou via provider
        return ConversationScreen(otherUserId: otherUserId);
      },
    ),
  ],
);
