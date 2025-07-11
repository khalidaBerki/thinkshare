import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/admin_dashboard_provider.dart';
import '../widgets/kpi_card.dart';
import '../widgets/bar_chart.dart';
import '../widgets/top_flop_list.dart';
import '../widgets/line_chart_media_progression.dart';
import 'package:intl/intl.dart';
import 'dart:math' as math;

class AdminDashboardScreen extends StatelessWidget {
  const AdminDashboardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => AdminDashboardProvider()..fetchStats(),
      child: Consumer<AdminDashboardProvider>(
        builder: (context, provider, _) {
          final colorScheme = Theme.of(context).colorScheme;
          return Scaffold(
            backgroundColor: colorScheme.background,
            appBar: AppBar(
              title: const Text(
                "Admin Dashboard",
                style: TextStyle(
                  fontFamily: 'Montserrat',
                  fontWeight: FontWeight.bold,
                ),
              ),
              backgroundColor: colorScheme.surface,
              elevation: 1,
              actions: [
                IconButton(
                  icon: const Icon(Icons.refresh),
                  tooltip: "Refresh",
                  onPressed: provider.fetchStats,
                ),
              ],
            ),
            body: provider.isLoading
                ? const Center(child: CircularProgressIndicator())
                : LayoutBuilder(
                    builder: (context, constraints) {
                      final isWide = constraints.maxWidth > 900;
                      return SingleChildScrollView(
                        padding: const EdgeInsets.all(28),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // Date range filter
                            Row(
                              children: [
                                Text(
                                  "Date range:",
                                  style: TextStyle(
                                    fontWeight: FontWeight.bold,
                                    fontSize: 16,
                                    fontFamily: 'Montserrat',
                                  ),
                                ),
                                const SizedBox(width: 12),
                                OutlinedButton.icon(
                                  icon: const Icon(Icons.date_range),
                                  label: Text(
                                    provider.selectedRange == null
                                        ? "All time"
                                        : "${provider.selectedRange!.start.toString().split(' ').first} - ${provider.selectedRange!.end.toString().split(' ').first}",
                                  ),
                                  onPressed: () async {
                                    final now = DateTime.now();
                                    final picked = await showDateRangePicker(
                                      context: context,
                                      firstDate: DateTime(now.year - 2),
                                      lastDate: now,
                                      initialDateRange: provider.selectedRange,
                                    );
                                    if (picked != null)
                                      provider.setRange(picked);
                                  },
                                ),
                              ],
                            ),
                            const SizedBox(height: 28),
                            // KPIs
                            Wrap(
                              spacing: 24,
                              runSpacing: 18,
                              children: [
                                KpiCard(
                                  label: "Total Posts",
                                  value: "${provider.totalPosts}",
                                  icon: Icons.article,
                                  color: Colors.green,
                                ),
                                KpiCard(
                                  label: "Total Comments",
                                  value: "${provider.totalComments}",
                                  icon: Icons.comment,
                                  color: Colors.deepPurple,
                                ),
                                KpiCard(
                                  label: "Total Conversations",
                                  value: "${provider.totalConversations}",
                                  icon: Icons.forum,
                                  color: Colors.blue,
                                ),
                              ],
                            ),
                            const SizedBox(height: 32),
                            // Graph
                            Text(
                              "Posts per day",
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 18,
                                fontFamily: 'Montserrat',
                              ),
                            ),
                            const SizedBox(height: 12),
                            BarChart(data: provider.postsPerDay),
                            const SizedBox(height: 32),
                            // Top/Flop
                            Flex(
                              direction: isWide
                                  ? Axis.horizontal
                                  : Axis.vertical,
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Expanded(
                                  child: TopFlopList(
                                    posts: provider.topPosts,
                                    label: "Top Posts",
                                    color: Colors.green,
                                  ),
                                ),
                                const SizedBox(width: 24, height: 24),
                                Expanded(
                                  child: TopFlopList(
                                    posts: provider.flopPosts,
                                    label: "Flop Posts",
                                    color: Colors.red,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 32),
                            // Content Preview
                            Text(
                              "Content Preview",
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 18,
                                fontFamily: 'Montserrat',
                              ),
                            ),
                            const SizedBox(height: 12),
                            ...provider.topPosts.map(
                              (p) => Card(
                                margin: const EdgeInsets.only(bottom: 16),
                                child: Padding(
                                  padding: const EdgeInsets.all(16),
                                  child: Row(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      // Title and date
                                      Expanded(
                                        flex: 2,
                                        child: Column(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              (p['document_type'] ??
                                                      p['documentType'] ??
                                                      (p['content']
                                                              ?.toString()
                                                              .substring(
                                                                0,
                                                                math.min(
                                                                  20,
                                                                  p['content']
                                                                          ?.toString()
                                                                          .length ??
                                                                      0,
                                                                ),
                                                              ) ??
                                                          '[No title]'))
                                                  .toString()
                                                  .toUpperCase(),
                                              style: const TextStyle(
                                                fontSize: 16,
                                                fontWeight: FontWeight.bold,
                                              ),
                                            ),
                                            const SizedBox(height: 4),
                                            Text(
                                              (p['created_at'] != null &&
                                                      p['created_at']
                                                          .toString()
                                                          .isNotEmpty)
                                                  ? DateFormat.yMMMd().format(
                                                      DateTime.tryParse(
                                                            p['created_at']
                                                                .toString(),
                                                          ) ??
                                                          DateTime(2000),
                                                    )
                                                  : '[No date]',
                                              style: const TextStyle(
                                                fontSize: 14,
                                                color: Colors.grey,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      const SizedBox(width: 16),
                                      // Content preview
                                      Expanded(
                                        child: Text(
                                          (p['content'] ?? '[No content]')
                                              .toString(),
                                          style: const TextStyle(fontSize: 14),
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                            const SizedBox(height: 32),
                            // Graphe courbe : progression des médias par jour
                            Text(
                              "Media progression (images/videos/docs)",
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 18,
                                fontFamily: 'Montserrat',
                              ),
                            ),
                            const SizedBox(height: 10),
                            LineChartMediaProgression(
                              data: provider.mediaProgression,
                            ),
                            const SizedBox(height: 32),
                          ],
                        ),
                      );
                    },
                  ),
          );
        },
      ),
    );
  }
}
