import 'package:app/utils/locales.dart';
import 'package:app/utils/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/vacancy_provider.dart';
import '../models/vacancy_model.dart';
import 'package:uuid/uuid.dart';

class CreateVacancyScreen extends StatefulWidget {
  const CreateVacancyScreen({super.key});

  @override
  CreateVacancyScreenState createState() => CreateVacancyScreenState();
}

class CreateVacancyScreenState extends State<CreateVacancyScreen> {
  final _formKey = GlobalKey<FormState>();
  final _descriptionController = TextEditingController();
  String _priceRange = '\$0 - \$50,000';
  int _yearsOfExperience = 0;
  String _workType = "Удаленная";

  @override
  Widget build(BuildContext context) {
    final vacancyProvider = Provider.of<VacancyProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocale.createVacancy.getString(context)),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _descriptionController,
                decoration: InputDecoration(labelText: AppLocale.vacancyDescription.getString(context)),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return AppLocale.enterDescription.getString(context);
                  }
                  return null;
                },
              ),
              DropdownButtonFormField<String>(
                value: _priceRange,
                decoration: InputDecoration(labelText: AppLocale.salaryRange.getString(context)),
                items: [
                  '\$0 - \$50,000',
                  '\$50,000 - \$100,000',
                  '\$100,000+',
                ].map((range) {
                  return DropdownMenuItem(value: range, child: Text(range));
                }).toList(),
                onChanged: (value) {
                  setState(() {
                    _priceRange = value!;
                  });
                },
              ),
              TextFormField(
                decoration: InputDecoration(labelText: AppLocale.yearsOfExperience.getString(context)),
                keyboardType: TextInputType.number,
                onChanged: (value) {
                  _yearsOfExperience = int.tryParse(value) ?? 0;
                },
              ),
              DropdownButtonFormField<String>(
                value: _workType,
                decoration: InputDecoration(labelText: AppLocale.workType.getString(context)),
                items: [AppLocale.remote.getString(context), AppLocale.office.getString(context)].map((type) {
                  return DropdownMenuItem(value: type, child: Text(type));
                }).toList(),
                onChanged: (value) {
                  setState(() {
                    _workType = value!;
                  });
                },
              ),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: () {
                  if (_formKey.currentState!.validate()) {
                    VacancyModel newVacancy = VacancyModel(
                      id: const Uuid().v4(),
                      description: _descriptionController.text,
                      priceRange: _priceRange,
                      yearsOfExperience: _yearsOfExperience,
                      workType: _workType,
                    );
                    vacancyProvider.createVacancy(newVacancy);
                    context.go(RouteNames.home);
                  }
                },
                child: Text(AppLocale.publishVacancy.getString(context)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
