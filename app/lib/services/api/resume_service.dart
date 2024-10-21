// services/api/resume_service.dart

import 'dart:convert';
import '../../models/resume.dart';
import 'api_client.dart';

class ResumeService {
  final ApiClient apiClient;

  ResumeService(this.apiClient);

  // Fetch all resumes
  Future<List<Resume>> getResumes() async {
    try {
      final response = await apiClient.get('resumes');

      if (response.statusCode == 200) {
        final List<dynamic> resumesJson = jsonDecode(response.body);

        return resumesJson.map((json) => Resume.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load resumes');
      }
    } catch (e) {
      throw e; // Rethrow the exception to be caught higher up
    }
  }

  // Fetch a single resume by ID
  Future<Resume> getResumeByID(int id) async {
    try {
      final response = await apiClient.get('resumes/$id');

      if (response.statusCode == 200) {
        //print(response.body);
        return Resume.fromJson(jsonDecode(response.body));
      } else {
        throw Exception('Failed to load resume');
      }
    } catch (e) {
      throw e;
    }
  }

  // Post a new resume
  Future<void> postResume(Map<String, dynamic> resume) async {
    try {
      final response = await apiClient.post('resumes', resume);
      if (response.statusCode != 200) {
        throw Exception('Failed to post resume');
      }
    } catch (e) {
      throw e;
    }
  }

  // You can add more methods like update or delete if needed.
}
