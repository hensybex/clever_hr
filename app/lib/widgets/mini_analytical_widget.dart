// lib/widgets/mini_analytical_widget.dart

import 'package:app/utils/locales.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import '../models/vacancy_resume_match.dart';
class MiniAnalyticalWidget extends StatefulWidget {
  final VacancyResumeMatch resumeMatch;
  const MiniAnalyticalWidget({super.key, required this.resumeMatch});
  @override
  MiniAnalyticalWidgetState createState() => MiniAnalyticalWidgetState();
}
class MiniAnalyticalWidgetState extends State<MiniAnalyticalWidget> {
  String? _hoveredMetric;
  int? _hoveredScore;
  @override
  Widget build(BuildContext context) {
    List<int> scores = [
      widget.resumeMatch.relevantWorkExperience.score,
      widget.resumeMatch.technicalSkillsAndProficiencies.score,
      widget.resumeMatch.educationAndCertifications.score,
      widget.resumeMatch.softSkillsAndCulturalFit.score,
      widget.resumeMatch.languageAndCommunicationSkills.score,
      widget.resumeMatch.problemSolvingAndAnalyticalAbilities.score,
      widget.resumeMatch.adaptabilityAndLearningCapacity.score,
      widget.resumeMatch.leadershipAndManagementExperience.score,
      widget.resumeMatch.motivationAndCareerObjectives.score,
      widget.resumeMatch.additionalQualificationsAndValueAdds.score,
    ];
    // Define the maximum possible score
    const int maxScore = 15;
    // Calculate the total score
    int totalScore = scores.reduce((a, b) => a + b);
    return LayoutBuilder(
      builder: (context, constraints) {
        double availableHeight = constraints.maxHeight;
        return SizedBox(
          height: availableHeight,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    // Bar graph with 10 bars
                    Expanded(
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: List.generate(scores.length, (index) {
                          double barHeight = (scores[index] / maxScore) * availableHeight;
                          return Flexible(
                            child: MouseRegion(
                              onEnter: (_) {
                                setState(() {
                                  _hoveredMetric = _getMetricNameByIndex(context, index); // Set hovered metric name
                                  _hoveredScore = scores[index]; // Set the score for the hovered metric
                                });
                              },
                              onExit: (_) {
                                setState(() {
                                  _hoveredMetric = null; // Clear hovered metric when mouse leaves
                                  _hoveredScore = null; // Clear the score
                                });
                              },
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.end, // Align bars at the bottom
                                children: [
                                  Container(
                                    margin: const EdgeInsets.symmetric(horizontal: 1),
                                    height: barHeight, // Set the calculated height based on the score
                                    color: Colors.black, // Customize the color as needed
                                  ),
                                ],
                              ),
                            ),
                          );
                        }),
                      ),
                    ),
                    // Display the total score next to the bar graph
                    Padding(
                      padding: const EdgeInsets.only(left: 10),
                      child: Text(
                        "${AppLocale.total.getString(context)} $totalScore",
                        style: const TextStyle(fontWeight: FontWeight.bold),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 5),
              // Display the hovered metric name and score below the bar graph
              Center(
                child: Text(
                  _hoveredMetric != null
                      ? "${AppLocale.metric.getString(context)} $_hoveredMetric $_hoveredScore" // Show metric name and score
                      : AppLocale.hoverOverTheBarForDetails.getString(context),
                  style: const TextStyle(fontSize: 12, fontStyle: FontStyle.italic),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
  String _getMetricNameByIndex(BuildContext context, int index) {
    switch (index) {
      case 0:
        return AppLocale.metricWorkExperience.getString(context);
      case 1:
        return AppLocale.metricTechnicalSkills.getString(context);
      case 2:
        return AppLocale.metricEducationCertifications.getString(context);
      case 3:
        return AppLocale.metricSoftSkills.getString(context);
      case 4:
        return AppLocale.metricLanguageSkills.getString(context);
      case 5:
        return AppLocale.metricProblemSolving.getString(context);
      case 6:
        return AppLocale.metricAdaptability.getString(context);
      case 7:
        return AppLocale.metricLeadershipExperience.getString(context);
      case 8:
        return AppLocale.metricMotivationObjectives.getString(context);
      case 9:
        return AppLocale.metricAdditionalQualifications.getString(context);
      default:
        return AppLocale.metricNumAnalysis.getString(context) + (index + 1).toString();
    }
  }
}
