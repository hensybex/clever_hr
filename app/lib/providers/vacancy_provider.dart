import 'package:app/services/hive/hive_init.dart';
import 'package:flutter/material.dart';
import '../models/vacancy_model.dart';
import '../services/api_service.dart';
import 'package:hive/hive.dart';

class VacancyProvider with ChangeNotifier {
  final ApiService _apiService = ApiService();
  List<VacancyModel> _vacancies = [];

  List<VacancyModel> get vacancies => _vacancies;

  Future<void> loadVacancies() async {
    // Load from Hive cache
    var box = Hive.box<VacancyModel>(BoxNames.vacancies);
    _vacancies = box.values.toList();

    // Fetch from API and update Hive
    var fetchedVacancies = await _apiService.fetchVacancies();
    _vacancies = fetchedVacancies;
    await box.clear();
    await box.addAll(_vacancies);

    notifyListeners();
  }

  Future<void> createVacancy(VacancyModel vacancy) async {
    var newVacancy = await _apiService.createVacancy(vacancy);
    _vacancies.add(newVacancy);

    // Update Hive cache
    var box = Hive.box<VacancyModel>(BoxNames.vacancies);
    await box.add(newVacancy);

    notifyListeners();

    // Establish WebSocket connection
    _apiService.publishVacancy(newVacancy.id, (message) {
      if (message['finished'] == true) {
        newVacancy.status = 'Analysed';
        box.put(newVacancy.key, newVacancy);
        notifyListeners();
      }
    });
  }
}
