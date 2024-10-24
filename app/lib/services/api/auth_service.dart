// lib/services/api/auth_service.dart

import 'dart:convert';
import 'package:http/http.dart' as http;
import '../hive/hive_init.dart';
import 'api_client.dart';
import 'package:hive_flutter/hive_flutter.dart';
import '../../models/user_token.dart';
class AuthService {
  final ApiClient apiClient;
  AuthService(this.apiClient);
  // Login method
  Future<Map<String, dynamic>> login(String username, String password) async {
    try {
      final Map<String, dynamic> body = {
        'username': username,
        'password': password,
      };
      // Send POST request to the login endpoint (no authentication needed)
      final http.Response response = await apiClient.post('login', body, needsAuth: false, isAuthEndpoint: true);
      if (response.statusCode == 200) {
        final Map<String, dynamic> jsonResponse = jsonDecode(response.body);
        // Assuming your API returns 'token' and 'expire' fields
        final String token = jsonResponse['token'];
        final DateTime expiry = DateTime.parse(jsonResponse['expire']);
        // Return the token and expiry date to the provider
        return {
          'token': token,
          'expire': expiry.toIso8601String(),
        };
      } else {
        throw Exception('Failed to login');
      }
    } catch (e) {
      throw Exception('Login failed');
    }
  }
  // Logout method
  Future<void> logout() async {
    try {
      // Send POST request to the logout endpoint (no authentication needed)
      final http.Response response = await apiClient.post('logout', {}, needsAuth: false, isAuthEndpoint: true);
      if (response.statusCode == 200) {
        // Remove token from Hive storage
        var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
        tokenBox.delete('token');
      } else {
        throw Exception('Failed to logout');
      }
    } catch (e) {
      throw Exception('Logout failed');
    }
  }
}
