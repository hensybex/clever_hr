// lib/screens/vacancy_resume_match_screen.dart

import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../models/vacancy_resume_match.dart';
import '../providers/analytical_widget_provider.dart';
import '../providers/match_provider.dart';
import '../utils/locales.dart';
import '../utils/router.dart';
import '../widgets/analytical_widget.dart';
import '../widgets/resume.dart';
class VacancyResumeMatchScreen extends StatefulWidget {
  final int matchId;
  const VacancyResumeMatchScreen({super.key, required this.matchId});
  @override
  VacancyResumeMatchScreenState createState() => VacancyResumeMatchScreenState();
}
class VacancyResumeMatchScreenState extends State<VacancyResumeMatchScreen> {
  @override
  void initState() {
    super.initState();
    _loadMatchAndResume();
  }
  Future<void> _loadMatchAndResume() async {
    final matchProvider = Provider.of<MatchProvider>(context, listen: false);
    await matchProvider.loadVacancyResumeMatch(widget.matchId);
    if (matchProvider.match != null) {
      await matchProvider.loadResume(matchProvider.match!.resumeId);
    }
  }
  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => AnalyticalWidgetProvider(), // Provide the hover state manager
      child: Scaffold(
        appBar: AppBar(
          title: Text(AppLocale.resumeAnalysis.getString(context)),
          leading: IconButton(
            icon: Icon(Icons.arrow_back),
            onPressed: () {
              final matchProvider = Provider.of<MatchProvider>(context, listen: false);
              final vacancyId = matchProvider.match?.vacancyId; // Get the vacancyId from matchProvider
              if (vacancyId != null) {
                // Use GoRouter to navigate back to the vacancy details screen
                WidgetsBinding.instance.addPostFrameCallback((_) {
                  context.go('${RouteNames.vacancy}/$vacancyId');
                });
              }
            },
          ),
        ),
        body: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Consumer<MatchProvider>(
            builder: (context, matchProvider, child) {
              final resumeMatch = matchProvider.match;
              if (resumeMatch == null) {
                return const Center(child: CircularProgressIndicator());
              }
              final resume = matchProvider.resume;
              if (resume == null) {
                return const Center(child: CircularProgressIndicator());
              }
              // Render the resume widget and the analytical widgets in rows
              return SingleChildScrollView(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ResumeWidget(resume: resume),
                    const SizedBox(height: 20),
                    _buildAnalyticalWidgetRows(context, resumeMatch),
                  ],
                ),
              );
            },
          ),
        ),
      ),
    );
  }
  // Build the rows with 2 AnalyticalWidgets per row
  Widget _buildAnalyticalWidgetRows(BuildContext context, VacancyResumeMatch resumeMatch) {
    List<Widget> rows = [];
    for (int i = 0; i < 10; i += 2) {
      rows.add(
        Row(
          children: [
            Expanded(child: _buildAnalyticalWidget(context, i, resumeMatch)),
            Expanded(child: _buildAnalyticalWidget(context, i + 1, resumeMatch)),
          ],
        ),
      );
    }
    return Column(children: rows);
  }
  // Helper method to build each AnalyticalWidget
  Widget _buildAnalyticalWidget(BuildContext context, int index, VacancyResumeMatch resumeMatch) {
    return Consumer<AnalyticalWidgetProvider>(
      builder: (context, provider, child) {
        bool isHovered = provider.getHoverState(index);
        int score = _getScoreByIndex(resumeMatch, index);
        String analysisText = _getAnalysisTextByIndex(context, index, resumeMatch);
        // Base height for non-hovered state
        double baseHeight = 150;
        // Estimate the number of lines required for the text
        int estimatedCharactersPerLine = 30; // Adjust this based on font size and widget width
        int totalLines = (analysisText.length / estimatedCharactersPerLine).ceil();
        // Calculate dynamic height when hovered, assuming each line takes about 20px
        double heightPerLine = 20.0;
        double calculatedHeight = baseHeight + (totalLines * heightPerLine);
        // Set a reasonable max height to avoid the widget becoming too large
        double maxHeight = 300; // Adjust this as needed for your layout
        double widgetHeight = isHovered ? calculatedHeight.clamp(baseHeight, maxHeight) : baseHeight;
        return AnalyticalWidget(
          score: score,
          analysisText: analysisText,
          isHovered: isHovered,
          onHoverChanged: (hovered) {
            provider.setHoverState(index, hovered);
          },
          widgetHeight: widgetHeight, // Pass calculated height to the widget
        );
      },
    );
  }
  int _getScoreByIndex(VacancyResumeMatch resumeMatch, int index) {
    switch (index) {
      case 0:
        return resumeMatch.relevantWorkExperience.score;
      case 1:
        return resumeMatch.technicalSkillsAndProficiencies.score;
      case 2:
        return resumeMatch.educationAndCertifications.score;
      case 3:
        return resumeMatch.softSkillsAndCulturalFit.score;
      case 4:
        return resumeMatch.languageAndCommunicationSkills.score;
      case 5:
        return resumeMatch.problemSolvingAndAnalyticalAbilities.score;
      case 6:
        return resumeMatch.adaptabilityAndLearningCapacity.score;
      case 7:
        return resumeMatch.leadershipAndManagementExperience.score;
      case 8:
        return resumeMatch.motivationAndCareerObjectives.score;
      case 9:
        return resumeMatch.additionalQualificationsAndValueAdds.score;
      default:
        return 0; // Provide a default score if index is out of range
    }
  }
  String _getAnalysisTextByIndex(BuildContext context, int index, VacancyResumeMatch resumeMatch) {
    switch (index) {
      case 0:
        return "${AppLocale.metricWorkExperience.getString(context)}\n${resumeMatch.relevantWorkExperience.overview}";
      case 1:
        return "${AppLocale.metricTechnicalSkills.getString(context)}\n${resumeMatch.technicalSkillsAndProficiencies.overview}";
      case 2:
        return "${AppLocale.metricEducationCertifications.getString(context)}\n${resumeMatch.educationAndCertifications.overview}";
      case 3:
        return "${AppLocale.metricSoftSkills.getString(context)}\n${resumeMatch.softSkillsAndCulturalFit.overview}";
      case 4:
        return "${AppLocale.metricLanguageSkills.getString(context)}\n${resumeMatch.languageAndCommunicationSkills.overview}";
      case 5:
        return "${AppLocale.metricProblemSolving.getString(context)}\n${resumeMatch.problemSolvingAndAnalyticalAbilities.overview}";
      case 6:
        return "${AppLocale.metricAdaptability.getString(context)}\n${resumeMatch.adaptabilityAndLearningCapacity.overview}";
      case 7:
        return "${AppLocale.metricLeadershipExperience.getString(context)}\n${resumeMatch.leadershipAndManagementExperience.overview}";
      case 8:
        return "${AppLocale.metricMotivationObjectives.getString(context)}\n${resumeMatch.motivationAndCareerObjectives.overview}";
      case 9:
        return "${AppLocale.metricAdditionalQualifications.getString(context)}\n${resumeMatch.additionalQualificationsAndValueAdds.overview}";
      default:
        return AppLocale.metricNumAnalysis.getString(context) + (index + 1).toString();
    }
  }
}
