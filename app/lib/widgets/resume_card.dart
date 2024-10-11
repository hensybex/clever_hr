import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import '../models/vacancy_resume_match_model.dart';
import '../utils/locales.dart';

class ResumeCard extends StatelessWidget {
  final VacancyResumeMatchModel resumeMatch;

  const ResumeCard({super.key, required this.resumeMatch});

  @override
  Widget build(BuildContext context) {
    // Calculate gradient based on analysedForMatch value
    double value = resumeMatch.analysedForMatch;
    Color startColor = Colors.blue;
    Color endColor = Colors.green;
    Color backgroundColor = Color.lerp(startColor, endColor, value)!;

    return Card(
      color: backgroundColor,
      child: ListTile(
        title: Text(AppLocale.resumeID.getString(context) + resumeMatch.resumeId.toString()),
        subtitle: Text(AppLocale.suitability.getString(context) + resumeMatch.suitability.toString()),
        onTap: () {
          context.go(RouteNames.vacancyResumeMatch, extra: resumeMatch);
        },
      ),
    );
  }
}
