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
  /* Future<Map<String, dynamic>> login(String username, String password) async {
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
  } */

  Future<Map<String, dynamic>> login(String username, String password) async {
  try {
    final Map<String, dynamic> body = {
      'username': username,
      'password': password,
    };

    // Log before sending the request
    print('Attempting to login with username: $username');

    // Send POST request to the login endpoint (no authentication needed)
    final http.Response response = await apiClient.post(
      'login',
      body,
      needsAuth: false,
      isAuthEndpoint: true,
    );

    // Check the status code and handle accordingly
    if (response.statusCode == 200) {
      final Map<String, dynamic> jsonResponse = jsonDecode(response.body);

      // Log successful login
      print('Login successful. Token received.');

      // Assuming your API returns 'token' and 'expire' fields
      final String token = jsonResponse['token'];
      final DateTime expiry = DateTime.parse(jsonResponse['expire']);

      // Save token to Hive storage
      var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
      tokenBox.put('token', UserToken(token: token, expiryDate: expiry));

      // Return the token and expiry date to the provider
      return {
        'token': token,
        'expire': expiry.toIso8601String(),
      };
    } else {
      // Log the error response
      print('Login failed with status code: ${response.statusCode}');
      print('Response body: ${response.body}');
      throw Exception('Failed to login');
    }
  } catch (e) {
    // Log the exception details
    print('Exception during login: $e');
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
