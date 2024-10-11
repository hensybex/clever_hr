import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import '../models/vacancy_model.dart';
import '../utils/locales.dart';

class VacancyCard extends StatelessWidget {
  final VacancyModel vacancy;

  const VacancyCard({super.key, required this.vacancy});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        title: Text(vacancy.description),
        subtitle: Text(AppLocale.status.getString(context) + vacancy.status.toString()),
        trailing: const Icon(Icons.arrow_forward),
        onTap: () {
          context.go(RouteNames.vacancy, extra: vacancy);
        },
      ),
    );
  }
}
