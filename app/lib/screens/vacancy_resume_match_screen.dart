import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import '../models/vacancy_resume_match_model.dart';
import '../utils/locales.dart';
import '../widgets/analytical_widget.dart';

class VacancyResumeMatchScreen extends StatelessWidget {
  final VacancyResumeMatchModel resumeMatch;

  const VacancyResumeMatchScreen({super.key, required this.resumeMatch});

  @override
  Widget build(BuildContext context) {
    List<Widget> analyticalWidgets = List.generate(10, (index) {
      return AnalyticalWidget(
        score: (resumeMatch.suitability * 100).toInt(),
        analysisText: AppLocale.metricNumAnalysis.getString(context) + (index + 1).toString(),
      );
    });

    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.resumeAnalysis.getString(context)),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Wrap(
          spacing: 10,
          runSpacing: 10,
          children: analyticalWidgets,
        ),
      ),
    );
  }
}
