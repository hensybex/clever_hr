import 'package:flutter/material.dart';

class AnalyticalWidget extends StatelessWidget {
  final int score; // From 1 to 100
  final String analysisText;

  const AnalyticalWidget({
    super.key,
    required this.score,
    required this.analysisText,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 150, // Adjust as needed
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            children: [
              LinearProgressIndicator(
                value: score / 100,
                minHeight: 8,
              ),
              const SizedBox(height: 8),
              Text(
                analysisText,
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
