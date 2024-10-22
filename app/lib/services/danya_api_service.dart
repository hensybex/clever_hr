// lib/services/danya_api_service.dart

import 'dart:convert';
import 'package:http/http.dart' as http;
import 'api/api_client.dart';

class DanyaApiService {
  final ApiClient _apiClient;

  DanyaApiService(this._apiClient);

  // Method to call the /llm/generate_interview_start endpoint
  Future<String> generateInterviewStart(int matchID, String tgLogin) async {
    final response = await _apiClient.post(
      'llm/generate_interview_start',
      {'match_id': matchID, 'tg_login': tgLogin},
      needsAuth: false, // Unprotected endpoint
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['message']; // Assuming the response contains a 'message' field
    } else {
      throw Exception('Failed to generate interview message');
    }
  }

  // Method to send a message to the external endpoint
  Future<void> sendMessage(String message) async {
    final url = 'https://hh.valiev.xyz/candidate/%40hensybex2/message/send';
    final response = await http.post(
      Uri.parse(url),
      headers: {
        'Content-Type': 'application/json',
        'accept': 'application/json',
      },
      body: jsonEncode({'message': message}),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to send message');
    }
  }

  // Combined method to generate and send the message
  Future<void> generateAndSendMessage(int matchID, String tgLogin) async {
    try {
      // Generate the interview message
      final generatedMessage = await generateInterviewStart(matchID, tgLogin);
      print('Generated message: $generatedMessage');

      // Send the generated message to the external endpoint
      await sendMessage(generatedMessage);
      print('Message sent successfully');
    } catch (error) {
      print('Error: $error');
      throw Exception('Failed to generate and send message');
    }
  }
}
