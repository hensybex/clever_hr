import 'dart:convert';

import '../models/vacancy_model.dart';
import '../models/vacancy_resume_match_model.dart';
import 'package:http/http.dart' as http;

class ApiService {
  final String apiBaseUrl;
  ApiService(this.apiBaseUrl);

  // Authentication
  Future<bool> login(String username, String password) async {
    return true;
  }

  // Vacancies
  Future<List<VacancyModel>> fetchVacancies() async {
    // TODO: Implement fetch vacancies API call
    return [];
  }

  Future<VacancyModel> createVacancy(VacancyModel vacancy) async {
    // TODO: Implement create vacancy API call
    return vacancy;
  }

  // WebSocket connection for publishing vacancy
  void publishVacancy(String vacancyId, Function onMessage) {
    // TODO: Implement WebSocket connection
  }

  // Resumes
  Future<List<VacancyResumeMatchModel>> fetchResumes(String vacancyId) async {
    // TODO: Implement fetch resumes API call
    return [];
  }
}
