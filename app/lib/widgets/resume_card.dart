// widgets/resume_card.dart

import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart'; // Assuming you're using go_router for navigation
import '../models/vacancy_resume_match.dart';
import '../utils/locales.dart';
import '../utils/router.dart';
import 'mini_analytical_widget.dart';

class ResumeCard extends StatelessWidget {
  final VacancyResumeMatch resumeMatch;

  const ResumeCard({super.key, required this.resumeMatch});

  @override
  Widget build(BuildContext context) {
    // Calculate gradient based on analysedForMatch value
    double value = resumeMatch.suitability;
    Color startColor = Colors.blue;
    Color endColor = Colors.green;
    Color backgroundColor = Color.lerp(startColor, endColor, value)!;

    return MouseRegion(
      cursor: SystemMouseCursors.click, // Update the cursor to pointer on hover
      child: GestureDetector(
        onTap: () {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            context.go('${RouteNames.vacancyResumeMatch}/${resumeMatch.id}');
          });
        },
        child: SizedBox(
          height: 150,
          child: Card(
            color: backgroundColor, // Keep the color static, no hover effect
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.stretch, // Stretch children to match card height
                children: [
                  // Title and Subtitle inside a container with fixed width
                  SizedBox(
                    width: 150, // Fixed width for the title and subtitle
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.center, // Center the title and subtitle vertically
                      children: [
                        Text(
                          AppLocale.resumeID.getString(context) + resumeMatch.resumeId.toString(),
                          style: const TextStyle(fontWeight: FontWeight.bold),
                        ),
                        Text(AppLocale.suitability.getString(context) + resumeMatch.suitability.toString()),
                      ],
                    ),
                  ),
                  const SizedBox(width: 10), // Small spacing between text and analytical widget
                  // Expanded to take all available space for MiniAnalyticalWidget
                  Expanded(
                    child: MiniAnalyticalWidget(resumeMatch: resumeMatch),
                  ),
                  // Add trailing arrow icon
                  const Icon(Icons.arrow_forward),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
