import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import '../models/vacancy_model.dart';
import '../models/vacancy_resume_match_model.dart';
import '../utils/locales.dart';
import '../widgets/resume_card.dart';

class VacancyScreen extends StatelessWidget {
  final VacancyModel vacancy;

  const VacancyScreen({super.key, required this.vacancy});

  @override
  Widget build(BuildContext context) {
    // Placeholder resumes list
    List<VacancyResumeMatchModel> resumes = []; // Fetch from provider or API

    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.vacancyDetails.getString(context)),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Text(
              vacancy.description,
              style: const TextStyle(fontSize: 18),
            ),
          ),
          Expanded(
            child: ListView.builder(
              itemCount: resumes.length,
              itemBuilder: (context, index) {
                return ResumeCard(resumeMatch: resumes[index]);
              },
            ),
          ),
        ],
      ),
    );
  }
}
