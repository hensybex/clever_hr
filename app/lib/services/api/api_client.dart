// lib/services/api/api_client.dart

import 'dart:convert';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:http/http.dart' as http;
import '../../models/user_token.dart';
import '../hive/hive_init.dart';

class ApiClient {
  final String baseUrl;
  ApiClient(this.baseUrl);

  // Get headers, but only when needed (for protected endpoints)
  Future<Map<String, String>> _getHeaders({bool needsAuth = true}) async {
    if (!needsAuth) {
      // Updated Content-Type header to match curl request
      return {'Content-Type': 'application/json'};
    }
    var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
    UserToken? userToken = tokenBox.get('token');
    if (userToken == null || userToken.token.isEmpty) {
      throw Exception("No valid token found.");
    }
    return {
      // Updated Content-Type header to match curl request
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ${userToken.token}',
    };
  }

  // Helper function to construct the full URL
  String _constructUrl(String path, {bool isAuthEndpoint = false}) {
    // If it's an auth endpoint, do not prepend the /api prefix
    if (isAuthEndpoint) {
      return '$baseUrl/$path';
    }
    // Otherwise, prepend the /api prefix to non-auth paths
    return '$baseUrl/api/$path';
  }

  // GET request
  Future<http.Response> get(String path, {bool needsAuth = true, bool isAuthEndpoint = false}) async {
    final headers = await _getHeaders(needsAuth: needsAuth);
    final url = _constructUrl(path, isAuthEndpoint: isAuthEndpoint);
    return await http.get(Uri.parse(url), headers: headers);
  }

  // POST request
  Future<http.Response> post(String path, Map<String, dynamic> body, 
    {bool needsAuth = true, bool isAuthEndpoint = false}) async {
    final headers = await _getHeaders(needsAuth: needsAuth);
    final url = _constructUrl(path, isAuthEndpoint: isAuthEndpoint);

    // Log the request details
    print('--- POST Request ---');
    print('URL: $url');
    print('Headers: $headers');
    print('Body: ${jsonEncode(body)}');

    final response = await http.post(Uri.parse(url), headers: headers, body: jsonEncode(body));

    // Log the response details
    print('--- Response ---');
    print('Status Code: ${response.statusCode}');
    print('Headers: ${response.headers}');
    print('Body: ${response.body}');

    return response;
  }

  // PUT request
  Future<http.Response> put(String path, Map<String, dynamic> body, 
    {bool needsAuth = true, bool isAuthEndpoint = false}) async {
    final headers = await _getHeaders(needsAuth: needsAuth);
    final url = _constructUrl(path, isAuthEndpoint: isAuthEndpoint);
    return await http.put(Uri.parse(url), headers: headers, body: jsonEncode(body));
  }
}
