import 'package:flutter/material.dart';

class AppTheme {
  static ThemeData lightTheme = ThemeData(
    fontFamily: 'Roboto',
    brightness: Brightness.light,
    primaryColor: Color(0xFF10A37F),
    scaffoldBackgroundColor: Colors.white,
    appBarTheme: const AppBarTheme(
      backgroundColor: Colors.white,
      foregroundColor: Colors.black,
      elevation: 0,
    ),
    colorScheme: const ColorScheme.light(
      primary: Color(0xFF10A37F),
      secondary: Color(0xFF10A37F),
    ),
    textTheme: Typography.blackCupertino,
  );

  static ThemeData darkTheme = ThemeData(
    fontFamily: 'Roboto',
    brightness: Brightness.dark,
    primaryColor: Color(0xFF10A37F),
    scaffoldBackgroundColor: Color(0xFF121212),
    appBarTheme: const AppBarTheme(
      backgroundColor: Color(0xFF121212),
      foregroundColor: Colors.white,
      elevation: 0,
    ),
    colorScheme: const ColorScheme.dark(
      primary: Color(0xFF10A37F),
      secondary: Color(0xFF10A37F),
    ),
    textTheme: Typography.whiteCupertino,
  );
}
