import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

class LineChartMediaProgression extends StatefulWidget {
  final Map<String, Map<String, int>>
  data; // {date: {image: x, video: y, doc: z}}
  const LineChartMediaProgression({super.key, required this.data});

  @override
  State<LineChartMediaProgression> createState() =>
      _LineChartMediaProgressionState();
}

class _LineChartMediaProgressionState extends State<LineChartMediaProgression>
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _animation;

  // Couleurs modernes avec dégradés
  final imageGradient = const LinearGradient(
    colors: [Color(0xFF2196F3), Color(0xFF0D47A1)],
  );

  final videoGradient = const LinearGradient(
    colors: [Color(0xFF9C27B0), Color(0xFF4A148C)],
  );

  final docGradient = const LinearGradient(
    colors: [Color(0xFFFF9800), Color(0xFFE65100)],
  );

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1500),
    );
    _animation = CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOutQuart,
    );
    _animationController.forward();
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final dates = widget.data.keys.toList()..sort();
    final imageSpots = <FlSpot>[];
    final videoSpots = <FlSpot>[];
    final docSpots = <FlSpot>[];

    // Prépare les données
    for (int i = 0; i < dates.length; i++) {
      final dayData = widget.data[dates[i]] ?? {};
      imageSpots.add(FlSpot(i.toDouble(), (dayData['image'] ?? 0).toDouble()));
      videoSpots.add(FlSpot(i.toDouble(), (dayData['video'] ?? 0).toDouble()));
      docSpots.add(FlSpot(i.toDouble(), (dayData['doc'] ?? 0).toDouble()));
    }

    return Column(
      children: [
        // Légende
        Padding(
          padding: const EdgeInsets.only(bottom: 20),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _legendItem("Images", const Color(0xFF2196F3)),
              const SizedBox(width: 24),
              _legendItem("Vidéos", const Color(0xFF9C27B0)),
              const SizedBox(width: 24),
              _legendItem("Documents", const Color(0xFFFF9800)),
            ],
          ),
        ),

        // Graphique
        SizedBox(
          height: 240,
          child: AnimatedBuilder(
            animation: _animation,
            builder: (context, _) {
              return Container(
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.surface,
                  borderRadius: BorderRadius.circular(16),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.grey.withOpacity(0.1),
                      blurRadius: 10,
                      spreadRadius: 1,
                    ),
                  ],
                ),
                padding: const EdgeInsets.only(top: 16, right: 16),
                child: LineChart(
                  LineChartData(
                    gridData: FlGridData(
                      show: true,
                      drawVerticalLine: true,
                      getDrawingHorizontalLine: (value) {
                        return FlLine(
                          color: Colors.grey.withOpacity(0.15),
                          strokeWidth: 0.8,
                          dashArray: [5, 5],
                        );
                      },
                      getDrawingVerticalLine: (value) {
                        return FlLine(
                          color: Colors.grey.withOpacity(0.15),
                          strokeWidth: 0.8,
                          dashArray: [5, 5],
                        );
                      },
                    ),
                    titlesData: FlTitlesData(
                      leftTitles: const AxisTitles(
                        sideTitles: SideTitles(
                          showTitles: true,
                          reservedSize: 36,
                          getTitlesWidget: _leftTitleWidgets,
                        ),
                      ),
                      bottomTitles: AxisTitles(
                        sideTitles: SideTitles(
                          showTitles: true,
                          getTitlesWidget: (value, meta) {
                            final idx = value.toInt();
                            if (idx < 0 || idx >= dates.length)
                              return const SizedBox();
                            final label = dates[idx].length >= 5
                                ? dates[idx].substring(5)
                                : dates[idx];
                            return Padding(
                              padding: const EdgeInsets.only(top: 8.0),
                              child: Text(
                                label,
                                style: TextStyle(
                                  fontSize: 11,
                                  fontWeight: FontWeight.w500,
                                  color: Colors.grey.shade700,
                                ),
                              ),
                            );
                          },
                          interval: dates.length > 7
                              ? (dates.length / 6).ceil().toDouble()
                              : 1,
                        ),
                      ),
                      rightTitles: const AxisTitles(
                        sideTitles: SideTitles(showTitles: false),
                      ),
                      topTitles: const AxisTitles(
                        sideTitles: SideTitles(showTitles: false),
                      ),
                    ),
                    borderData: FlBorderData(show: false),
                    minX: 0,
                    maxX: (dates.length - 1).toDouble(),
                    lineTouchData: LineTouchData(
                      touchTooltipData: LineTouchTooltipData(
                        tooltipBgColor: Colors.black.withOpacity(0.8),
                        tooltipRoundedRadius: 8,
                        tooltipPadding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                        getTooltipItems: (touchedSpots) {
                          return touchedSpots.map((touchedSpot) {
                            final String typeText;
                            switch (touchedSpot.barIndex) {
                              case 0:
                                typeText = "Images";
                                break;
                              case 1:
                                typeText = "Vidéos";
                                break;
                              case 2:
                                typeText = "Docs";
                                break;
                              default:
                                typeText = "";
                            }
                            return LineTooltipItem(
                              "$typeText: ${touchedSpot.y.toInt()}",
                              const TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                              ),
                            );
                          }).toList();
                        },
                      ),
                    ),
                    lineBarsData: [
                      _createLineData(imageSpots, imageGradient.colors, 0),
                      _createLineData(videoSpots, videoGradient.colors, 1),
                      _createLineData(docSpots, docGradient.colors, 2),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  LineChartBarData _createLineData(
    List<FlSpot> spots,
    List<Color> colors,
    int index,
  ) {
    // Animation de progression
    final animatedSpots = spots.asMap().entries.map((entry) {
      final index = entry.key;
      final spot = entry.value;
      final double animValue = _animation.value;
      if (index <= spots.length * animValue) {
        return spot;
      } else {
        return FlSpot(spot.x, 0);
      }
    }).toList();

    return LineChartBarData(
      spots: animatedSpots,
      isCurved: true,
      curveSmoothness: 0.35, // Courbes plus fluides et modernes
      gradient: LinearGradient(colors: colors),
      barWidth: 3.5,
      isStrokeCapRound: true,
      dotData: FlDotData(
        show: true,
        getDotPainter: (spot, percent, barData, index) {
          return FlDotCirclePainter(
            radius: 3,
            color: colors[0],
            strokeWidth: 1,
            strokeColor: Colors.white,
          );
        },
        checkToShowDot: (spot, barData) {
          // Ne montre que quelques points pour ne pas surcharger
          return spot.x % 2 == 0 || spot.x == barData.spots.last.x;
        },
      ),
      belowBarData: BarAreaData(
        show: true,
        gradient: LinearGradient(
          colors: [colors[0].withOpacity(0.2), colors[1].withOpacity(0.05)],
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
        ),
      ),
    );
  }

  Widget _legendItem(String label, Color color) {
    return Row(
      children: [
        Container(
          width: 12,
          height: 12,
          decoration: BoxDecoration(color: color, shape: BoxShape.circle),
        ),
        const SizedBox(width: 6),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w500,
            color: Colors.grey.shade800,
          ),
        ),
      ],
    );
  }
}

Widget _leftTitleWidgets(double value, TitleMeta meta) {
  final style = TextStyle(
    fontSize: 11,
    fontWeight: FontWeight.w500,
    color: Colors.grey.shade700,
  );

  return Text(
    value.toInt().toString(),
    style: style,
    textAlign: TextAlign.center,
  );
}
