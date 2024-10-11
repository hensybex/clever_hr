import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/vacancy_provider.dart';
import '../utils/locales.dart';
import '../widgets/vacancy_card.dart';
import '../models/vacancy_model.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final vacancyProvider = Provider.of<VacancyProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.vacancies.getString(context)),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              context.go(RouteNames.createVacancy);
            },
          ),
        ],
      ),
      body: FutureBuilder(
        future: vacancyProvider.loadVacancies(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            return ListView.builder(
              itemCount: vacancyProvider.vacancies.length,
              itemBuilder: (context, index) {
                VacancyModel vacancy = vacancyProvider.vacancies[index];
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
