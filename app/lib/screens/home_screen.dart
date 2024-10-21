// lib/screens/home_screen.dart

import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../providers/localization_provider.dart';
import '../providers/vacancy_provider.dart';
import '../utils/locales.dart';
import '../widgets/vacancy_card.dart';
import '../models/vacancy.dart';
class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});
  @override
  HomeScreenState createState() => HomeScreenState();
}
class HomeScreenState extends State<HomeScreen> {
  late Future<void> _vacancyFuture;
  @override
  void initState() {
    super.initState();
    final vacancyProvider = Provider.of<VacancyProvider>(context, listen: false);
    _vacancyFuture = vacancyProvider.loadVacancies();
  }
  @override
  Widget build(BuildContext context) {
    final vacancyProvider = Provider.of<VacancyProvider>(context);
    final authProvider = Provider.of<AuthProvider>(context, listen: false);
    final localizationProvider = Provider.of<LocalizationProvider>(context); // Access the localization provider
    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.vacancies.getString(context)), // Use localized text
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            color: Colors.white,
            onPressed: () {
              context.go(RouteNames.createVacancy);
            },
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            color: Colors.redAccent,
            onPressed: () {
              authProvider.logout();
              context.go(RouteNames.login);
            },
          ),
          ElevatedButton(
            onPressed: () {
              // Switch between 'en' and 'ru' based on the current language
              final newLanguageCode = localizationProvider.currentLanguage == 'ru' ? 'en' : 'ru';
              localizationProvider.switchLanguage(newLanguageCode);
            },
            child: Text(localizationProvider.currentLanguageName), // Display the current language name
          ),
        ],
      ),
      body: FutureBuilder(
        future: _vacancyFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            if (vacancyProvider.vacancies.isEmpty) {
              return Center(
                  //child: Text(AppLocale.noVacancies.getString(context)), // Use localized text
                  );
            }
            return ListView.builder(
              itemCount: vacancyProvider.vacancies.length,
              itemBuilder: (context, index) {
                Vacancy vacancy = vacancyProvider.vacancies[index];
                return VacancyCard(vacancy: vacancy);
              },
            );
          }
          return const Center(child: CircularProgressIndicator());
        },
      ),
    );
  }
}
