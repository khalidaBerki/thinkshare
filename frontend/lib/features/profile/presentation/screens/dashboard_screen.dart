import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/profile_provider.dart';
import 'package:intl/intl.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  DateTimeRange? _selectedRange;

  @override
  void initState() {
    super.initState();
    // On pourrait charger les posts ici si besoin
    final provider = Provider.of<ProfileProvider>(context, listen: false);
    if (provider.myPosts.isEmpty) {
      provider.fetchMyPosts();
    }
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final provider = Provider.of<ProfileProvider>(context);

    // Filtrage des posts par date si range sélectionné
    List<Map<String, dynamic>> posts = provider.myPosts;
    if (_selectedRange != null) {
      posts = posts.where((post) {
        final createdAt = DateTime.tryParse(post['created_at'] ?? '');
        if (createdAt == null) return false;
        return createdAt.isAfter(
              _selectedRange!.start.subtract(const Duration(days: 1)),
            ) &&
            createdAt.isBefore(
              _selectedRange!.end.add(const Duration(days: 1)),
            );
      }).toList();
    }

    // KPIs
    final totalPosts = posts.length;
    final followers = provider.myProfile?['followers'] ?? 0;
    final totalLikes = posts.fold<int>(
      0,
      (sum, p) => sum + ((p['like_count'] ?? 0) as int),
    );
    final totalComments = posts.fold<int>(
      0,
      (sum, p) => sum + ((p['comment_count'] ?? 0) as int),
    );
    final avgEngagement = totalPosts > 0
        ? ((totalLikes + totalComments) / totalPosts).toStringAsFixed(1)
        : '0';

    // Top/Flop posts
    final topPost = posts.isNotEmpty
        ? posts.reduce(
            (a, b) =>
                ((a['like_count'] ?? 0) + (a['comment_count'] ?? 0)) >=
                    ((b['like_count'] ?? 0) + (b['comment_count'] ?? 0))
                ? a
                : b,
          )
        : null;
    final flopPost = posts.isNotEmpty
        ? posts.reduce(
            (a, b) =>
                ((a['like_count'] ?? 0) + (a['comment_count'] ?? 0)) <=
                    ((b['like_count'] ?? 0) + (b['comment_count'] ?? 0))
                ? a
                : b,
          )
        : null;

    // Breakdown
    int images = 0, videos = 0, texts = 0;
    for (final p in posts) {
      if ((p['media_urls'] as List?)?.isNotEmpty == true) {
        // Fake: si url contient .mp4 ou .mov => vidéo, sinon image
        final urls = List<String>.from(p['media_urls']);
        if (urls.any((u) => u.endsWith('.mp4') || u.endsWith('.mov'))) {
          videos++;
        } else {
          images++;
        }
      } else {
        texts++;
      }
    }
    final totalContent = images + videos + texts;
    final imgPct = totalContent > 0 ? (images * 100 ~/ totalContent) : 0;
    final vidPct = totalContent > 0 ? (videos * 100 ~/ totalContent) : 0;
    final txtPct = totalContent > 0 ? (texts * 100 ~/ totalContent) : 0;

    // Graph data (posts per day)
    final Map<String, int> postsPerDay = {};
    for (final p in posts) {
      final date = DateTime.tryParse(p['created_at'] ?? '');
      if (date != null) {
        final key = DateFormat('yyyy-MM-dd').format(date);
        postsPerDay[key] = (postsPerDay[key] ?? 0) + 1;
      }
    }

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(
            Icons.arrow_back_ios_new_rounded,
            color: Colors.black87,
          ),
          onPressed: () => context.go('/profile'),
          tooltip: "Back",
        ),
        title: const Text(
          "Dashboard",
          style: TextStyle(
            fontWeight: FontWeight.bold,
            fontFamily: 'Montserrat',
            color: Colors.black87,
            fontSize: 22,
            letterSpacing: 0.5,
          ),
        ),
        centerTitle: true,
        backgroundColor: colorScheme.surface,
        elevation: 1,
        shadowColor: colorScheme.primary.withOpacity(0.06),
        surfaceTintColor: colorScheme.primary,
        actions: [
          IconButton(
            icon: const Icon(Icons.date_range),
            tooltip: "Filter by date",
            onPressed: () async {
              final now = DateTime.now();
              final picked = await showDateRangePicker(
                context: context,
                firstDate: DateTime(now.year - 2),
                lastDate: now,
                initialDateRange: _selectedRange,
              );
              if (picked != null) setState(() => _selectedRange = picked);
            },
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(18),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Card(
              elevation: 2,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16),
              ),
              color: colorScheme.surfaceVariant,
              child: Padding(
                padding: const EdgeInsets.symmetric(
                  vertical: 18,
                  horizontal: 16,
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      "Key Statistics",
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 18,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(
                          child: _statCard(
                            context,
                            "Total posts",
                            "$totalPosts",
                            Icons.article_outlined,
                          ),
                        ),
                        Expanded(
                          child: _statCard(
                            context,
                            "Followers",
                            "$followers",
                            Icons.people,
                          ),
                        ),
                        Expanded(
                          child: _statCard(
                            context,
                            "Avg. engagement",
                            "$avgEngagement",
                            Icons.trending_up,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 22),
            const Text(
              "Posts per day",
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const SizedBox(height: 10),
            _simpleBarChart(postsPerDay, colorScheme),
            const SizedBox(height: 22),
            const Text(
              "Content breakdown",
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const SizedBox(height: 10),
            Row(
              children: [
                _breakdownBadge(context, "$imgPct% images", Colors.blue),
                const SizedBox(width: 8),
                _breakdownBadge(context, "$vidPct% videos", Colors.deepPurple),
                const SizedBox(width: 8),
                _breakdownBadge(context, "$txtPct% texts", Colors.orange),
              ],
            ),
            const SizedBox(height: 22),
            const Text(
              "Top & Flop Posts",
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const SizedBox(height: 10),
            Row(
              children: [
                Expanded(
                  child: topPost != null
                      ? _postSummaryCard(context, topPost, "Top post")
                      : _emptyCard("No posts"),
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: flopPost != null
                      ? _postSummaryCard(context, flopPost, "Flop post")
                      : _emptyCard("No posts"),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _statCard(
    BuildContext context,
    String label,
    String value,
    IconData icon,
  ) {
    final colorScheme = Theme.of(context).colorScheme;
    return Card(
      elevation: 0,
      color: colorScheme.background,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 6),
        child: Column(
          children: [
            Icon(icon, color: colorScheme.primary, size: 28),
            const SizedBox(height: 6),
            Text(
              value,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 18,
                color: colorScheme.primary,
              ),
            ),
            const SizedBox(height: 2),
            Text(
              label,
              style: const TextStyle(fontSize: 13, color: Colors.grey),
            ),
          ],
        ),
      ),
    );
  }

  Widget _breakdownBadge(BuildContext context, String text, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.12),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        text,
        style: TextStyle(
          color: color,
          fontWeight: FontWeight.w600,
          fontSize: 13,
        ),
      ),
    );
  }

  Widget _postSummaryCard(
    BuildContext context,
    Map<String, dynamic> post,
    String label,
  ) {
    final colorScheme = Theme.of(context).colorScheme;
    final likeCount = post['like_count'] ?? 0;
    final commentCount = post['comment_count'] ?? 0;
    final content = post['content'] ?? '';
    final createdAt = post['created_at'] ?? '';
    return Card(
      color: colorScheme.surfaceVariant,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              label,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: colorScheme.primary,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              content.length > 60 ? "${content.substring(0, 60)}..." : content,
              style: const TextStyle(fontSize: 13),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Icon(Icons.star, color: colorScheme.primary, size: 18),
                const SizedBox(width: 4),
                Text("$likeCount"),
                const SizedBox(width: 12),
                Icon(
                  Icons.mode_comment_outlined,
                  color: colorScheme.primary,
                  size: 18,
                ),
                const SizedBox(width: 4),
                Text("$commentCount"),
                const SizedBox(width: 12),
                Icon(
                  Icons.calendar_today,
                  color: colorScheme.primary,
                  size: 16,
                ),
                const SizedBox(width: 4),
                Text(createdAt.toString().split('T').first),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _emptyCard(String text) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(18),
        child: Center(
          child: Text(text, style: const TextStyle(color: Colors.grey)),
        ),
      ),
    );
  }

  Widget _simpleBarChart(Map<String, int> data, ColorScheme colorScheme) {
    if (data.isEmpty) {
      return Container(
        height: 80,
        decoration: BoxDecoration(
          color: colorScheme.surfaceVariant,
          borderRadius: BorderRadius.circular(16),
        ),
        child: const Center(child: Text("No data")),
      );
    }
    final maxVal = data.values.fold<int>(0, (max, v) => v > max ? v : max);
    return Container(
      height: 130, // Augmente la hauteur ici
      padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
      decoration: BoxDecoration(
        color: colorScheme.surfaceVariant,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: data.entries.map((e) {
          final pct = maxVal > 0 ? e.value / maxVal : 0.0;
          return Expanded(
            child: Padding(
              padding: const EdgeInsets.only(
                bottom: 4,
              ), // Réduit le padding ici
              child: Column(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Container(
                    height: 80 * pct,
                    width: 12,
                    decoration: BoxDecoration(
                      color: colorScheme.primary,
                      borderRadius: BorderRadius.circular(6),
                    ),
                  ),
                  Text(
                    e.key.substring(5),
                    style: const TextStyle(fontSize: 10),
                  ),
                ],
              ),
            ),
          );
        }).toList(),
      ),
    );
  }
}
