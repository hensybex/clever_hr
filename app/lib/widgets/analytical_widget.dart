// widgets/analytical_widget.dart

import 'package:flutter/material.dart';

class AnalyticalWidget extends StatelessWidget {
  final int score;
  final String analysisText;
  final bool isHovered; // Hover state from provider
  final Function(bool) onHoverChanged;
  final double widgetHeight; // Pass dynamic height to widget

  const AnalyticalWidget({
    super.key,
    required this.score,
    required this.analysisText,
    required this.isHovered,
    required this.onHoverChanged,
    required this.widgetHeight, // Use dynamic height
  });

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => onHoverChanged(true),
      onExit: (_) => onHoverChanged(false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        width: 150,
        height: widgetHeight, // Use dynamic height based on hover state
        child: Card(
          child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Column(
              children: [
                LinearProgressIndicator(
                  value: score / 10,
                  minHeight: 8,
                ),
                const SizedBox(height: 8),
                Expanded(
                  child: Text(
                    analysisText + "\n\n" + "Оценка: " + score.toString(),
                    textAlign: TextAlign.center,
                    // Use ellipsis when not hovered, show full text when hovered
                    maxLines: isHovered ? null : 3, // Limit lines when not hovered
                    overflow: isHovered ? TextOverflow.visible : TextOverflow.ellipsis, // Ellipsis when not hovered
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
