// lib/screens/create_vacancy_screen.dart

import 'package:app/utils/locales.dart';
import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/vacancy_provider.dart';
import '../models/vacancy.dart';

class CreateVacancyScreen extends StatefulWidget {
  const CreateVacancyScreen({super.key});
  @override
  CreateVacancyScreenState createState() => CreateVacancyScreenState();
}
class CreateVacancyScreenState extends State<CreateVacancyScreen> {
  final _formKey = GlobalKey<FormState>();
  final _descriptionController = TextEditingController();
  @override
  Widget build(BuildContext context) {
    final vacancyProvider = Provider.of<VacancyProvider>(context);
    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.createVacancy.getString(context)),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            // Show a loading indicator if the provider is loading
            if (vacancyProvider.isLoading) const LinearProgressIndicator(),
            const SizedBox(height: 16),
            Expanded(
              child: Form(
                key: _formKey,
                child: Column(
                  children: [
                    // Description Text Field
                    TextFormField(
                      controller: _descriptionController,
                      maxLines: 10,
                      decoration: InputDecoration(
                        labelText: AppLocale.vacancyDescription.getString(context),
                        border: const OutlineInputBorder(),
                        alignLabelWithHint: true,
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return AppLocale.enterDescription.getString(context);
                        }
                        return null;
                      },
                    ),
                    const SizedBox(height: 20),
                    // Display any error messages
                    if (vacancyProvider.errorMessage.isNotEmpty)
                      Text(
                        vacancyProvider.errorMessage,
                        style: const TextStyle(color: Colors.red),
                      ),
                    // Display waiting messages or results
                    if (vacancyProvider.waitingMessage.isNotEmpty) Text(vacancyProvider.waitingMessage),
                    if (vacancyProvider.displayedContent.isNotEmpty) Text(vacancyProvider.displayedContent),
                    const SizedBox(height: 20),
                    // Publish Button
                    ElevatedButton(
                      onPressed: () {
                        if (_formKey.currentState!.validate()) {
                          Vacancy newVacancy = Vacancy(
                            description: _descriptionController.text,
                          );
                          vacancyProvider.createVacancy(newVacancy);
                          // Navigate to home after vacancy creation
                          context.go(RouteNames.home);
                        }
                      },
                      child: Text(AppLocale.publishVacancy.getString(context)),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
