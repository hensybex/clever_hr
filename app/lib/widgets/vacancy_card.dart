// lib/widgets/vacancy_card.dart

import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import '../models/vacancy.dart';
import '../utils/locales.dart';
class VacancyCard extends StatelessWidget {
  final Vacancy vacancy;
  const VacancyCard({super.key, required this.vacancy});
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 100,
      width: 300,
      child: Card(
        child: ListTile(
          title: Text(
            vacancy.title ?? 'No Title',
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          subtitle: Text(
            AppLocale.status.getString(context) + vacancy.status.toString(),
          ),
          trailing: const Icon(Icons.arrow_forward),
          onTap: () {
            WidgetsBinding.instance.addPostFrameCallback((_) {
              context.go('${RouteNames.vacancy}/${vacancy.id}', extra: vacancy);
            });
          },
        ),
      ),
    );
  }
}
