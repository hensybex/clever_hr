// lib/screens/vacancy_screen.dart

import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/vacancy_provider.dart';
import '../utils/locales.dart';
import '../utils/router.dart';
import '../widgets/expandable_text.dart';
import '../widgets/resume_card.dart';
class VacancyScreen extends StatefulWidget {
  final int vacancyId;
  const VacancyScreen({super.key, required this.vacancyId});
  @override
  VacancyScreenState createState() => VacancyScreenState();
}
class VacancyScreenState extends State<VacancyScreen> {
  bool isExpanded = false;
  late int vacancyId;
  @override
  void initState() {
    super.initState();
    vacancyId = widget.vacancyId;
    final vacancyProvider = Provider.of<VacancyProvider>(context, listen: false);
    vacancyProvider.fetchVacancyResumeMatches(vacancyId);
    vacancyProvider.fetchVacancy(vacancyId);
  }
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.vacancyScreen.getString(context)),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back), // Icon for the back button
          onPressed: () {
            context.go(RouteNames.home); // Navigate to the home route
          },
        ),
      ),
      body: Consumer<VacancyProvider>(
        builder: (context, vacancyProvider, child) {
          if (vacancyProvider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (vacancyProvider.errorMessage.isNotEmpty) {
            return Center(child: Text(vacancyProvider.errorMessage));
          }
          final matches = vacancyProvider.vacancyResumeMatches;
          return SingleChildScrollView(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    AppLocale.vacancyName.getString(context),
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                  GestureDetector(
                    onTap: () {
                      setState(() {
                        isExpanded = !isExpanded;
                      });
                    },
                    child: Text(
                      vacancyProvider.vacancy!.title!,
                      maxLines: isExpanded ? null : 1,
                      overflow: isExpanded ? TextOverflow.visible : TextOverflow.ellipsis,
                    ),
                  ),
                  ExpandableTextWidget(
                    label: AppLocale.vacancyDescription.getString(context),
                    text: vacancyProvider.vacancy!.description,
                  ),
                  // Wrap ListView.builder inside a Container or set shrinkWrap and physics
                  ListView.builder(
                    shrinkWrap: true, // Allow the ListView to shrink its height
                    physics: const NeverScrollableScrollPhysics(), // Disable scrolling for the ListView
                    itemCount: matches.length,
                    itemBuilder: (context, index) {
                      return ResumeCard(resumeMatch: matches[index]);
                    },
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
