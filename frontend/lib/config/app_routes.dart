import 'package:go_router/go_router.dart';

import '../features/auth/presentation/screens/entry_screen.dart';
//import '../features/auth/presentation/screens/login_screen.dart';
//import '../features/auth/presentation/screens/register_screen.dart';

final GoRouter appRouter = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(path: '/', builder: (context, state) => const EntryScreen()),

    /* GoRoute(
      path: '/login', 
      builder: (context, state) => const LoginScreen()
    ),
    
    
    GoRoute(
      path: '/register',
      builder: (context, state) => const RegisterScreen(),
    ),*/
  ],
);
