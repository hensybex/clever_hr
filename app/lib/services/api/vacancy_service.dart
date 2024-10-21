// services/api/vacancy_service.dart

import 'dart:convert';
import '../../models/vacancy.dart';
import 'api_client.dart';

class VacancyService {
  final ApiClient apiClient;

  VacancyService(this.apiClient);

  Future<List<Vacancy>> getVacancies() async {
    try {
      final response = await apiClient.get('vacancies');

      if (response.statusCode == 200) {
        final List<dynamic> vacanciesJson = jsonDecode(response.body);

        return vacanciesJson.map((json) => Vacancy.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load vacancies');
      }
    } catch (e) {
      throw e; // Rethrow the exception to be caught higher up
    }
  }

  Future<void> postVacancy(Map<String, dynamic> vacancy) async {
    final response = await apiClient.post('vacancies', vacancy);
    if (response.statusCode != 200) {
      throw Exception('Failed to post vacancy');
    }
  }

  Future<Vacancy> getVacancyByID(int id) async {
    final response = await apiClient.get('vacancies/$id');
    if (response.statusCode == 200) {
      return Vacancy.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load vacancy');
    }
  }
}
