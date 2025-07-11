import 'package:go_router/go_router.dart';

import '../features/auth/presentation/screens/entry_screen.dart';
import '../features/auth/presentation/screens/login_screen.dart';
import '../features/auth/presentation/screens/register_screen.dart';
import '../features/home/presentation/screens/post_detail_screen.dart';
import '../features/profile/presentation/screens/user_profile_screen.dart';
import '../../features/profile/presentation/screens/dashboard_screen.dart';
import '../../features/admin_dashboard/presentation/screens/admin_dashboard_screen.dart';

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
      builder: (context, state) => const NavigationScreen(tabIndex: 0),
    ),

    GoRoute(
      path: '/post',
      builder: (context, state) => const NavigationScreen(tabIndex: 1),
    ),

    GoRoute(
      path: '/messages',
      builder: (context, state) => const NavigationScreen(tabIndex: 2),
    ),

    GoRoute(
      path: '/profile',
      builder: (context, state) => const NavigationScreen(tabIndex: 3),
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
        return ConversationScreen(otherUserId: otherUserId);
      },
    ),
    GoRoute(
      path: '/user/:id',
      builder: (context, state) {
        final id = int.parse(state.pathParameters['id']!);
        return UserProfileScreen(userId: id);
      },
    ),
    GoRoute(
      path: '/dashboard',
      builder: (context, state) => const DashboardScreen(),
    ),
    GoRoute(
      path: '/admin',
      builder: (context, state) => const AdminDashboardScreen(),
    ),
  ],
);
