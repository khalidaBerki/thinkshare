import 'package:flutter/material.dart';

class TopFlopList extends StatelessWidget {
  final List<Map<String, dynamic>> posts;
  final String label;
  final Color color;
  const TopFlopList({
    super.key,
    required this.posts,
    required this.label,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 0,
      color: color.withOpacity(0.06),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              label,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: color,
                fontSize: 16,
              ),
            ),
            const SizedBox(height: 10),
            ...posts.map(
              (p) => Padding(
                padding: const EdgeInsets.symmetric(vertical: 4),
                child: Row(
                  children: [
                    Expanded(
                      child: Text(
                        (p['content'] ?? '[No content]').toString(),
                        style: const TextStyle(fontSize: 14),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    Icon(Icons.star, color: color, size: 18),
                    Text(" ${p['like_count'] ?? 0}"),
                    const SizedBox(width: 10),
                    Icon(Icons.comment, color: color, size: 18),
                    Text(" ${p['comment_count'] ?? 0}"),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
