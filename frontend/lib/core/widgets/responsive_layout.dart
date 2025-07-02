import 'package:flutter/material.dart';

class ResponsiveLayout extends StatelessWidget {
  final Widget mobileBody;
  final Widget webBody;

  const ResponsiveLayout({
    super.key,
    required this.mobileBody,
    required this.webBody,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        // Exemple seuil arbitraire (tu peux l'ajuster)
        if (constraints.maxWidth < 800) {
          return mobileBody;
        } else {
          return webBody;
        }
      },
    );
  }
}
