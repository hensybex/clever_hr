// services/api/match_service.dart

import 'dart:convert';
import '../../models/vacancy_resume_match.dart';
import 'api_client.dart';

class MatchService {
  final ApiClient apiClient;

  MatchService(this.apiClient);

  // Fetch resume matches for a specific vacancy
  Future<List<VacancyResumeMatch>> fetchVacancyResumeMatches(int vacancyId) async {
    try {
      final response = await apiClient.get('match/$vacancyId/matches');
      if (response.statusCode == 200) {
        final List<dynamic> matchesJson = jsonDecode(response.body);

        return matchesJson.map((json) => VacancyResumeMatch.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load resume matches');
      }
    } catch (e) {
      throw Exception('Failed to load resume matches');
    }
  }

  // Fetch details of a specific resume match by matchId
  Future<VacancyResumeMatch> fetchVacancyResumeMatchById(int matchId) async {
    try {
      final response = await apiClient.get('match/details/$matchId');

      if (response.statusCode == 200) {
        final Map<String, dynamic> matchJson = jsonDecode(response.body);

        return VacancyResumeMatch.fromJson(matchJson);
      } else {
        throw Exception('Failed to load resume match details');
      }
    } catch (e) {
      throw Exception('Failed to load resume match details');
    }
  }
}
