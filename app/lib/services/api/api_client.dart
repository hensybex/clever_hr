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
      return {'Content-Type': 'application/json; charset=UTF-8'};
    }
    var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
    UserToken? userToken = tokenBox.get('token');
    if (userToken == null || userToken.token.isEmpty) {
      throw Exception("No valid token found.");
    }
    return {
      'Content-Type': 'application/json; charset=UTF-8',
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
  Future<http.Response> get(String path, {bool needsAuth = true, bool isAuthEndpoint = false}) async {
    final headers = await _getHeaders(needsAuth: needsAuth);
    final url = _constructUrl(path, isAuthEndpoint: isAuthEndpoint);
    return await http.get(Uri.parse(url), headers: headers);
  }
  Future<http.Response> post(String path, Map<String, dynamic> body, {bool needsAuth = true, bool isAuthEndpoint = false}) async {
    final headers = await _getHeaders(needsAuth: needsAuth);
    final url = _constructUrl(path, isAuthEndpoint: isAuthEndpoint);
    return await http.post(Uri.parse(url), headers: headers, body: jsonEncode(body));
  }
  Future<http.Response> put(String path, Map<String, dynamic> body, {bool needsAuth = true, bool isAuthEndpoint = false}) async {
    final headers = await _getHeaders(needsAuth: needsAuth);
    final url = _constructUrl(path, isAuthEndpoint: isAuthEndpoint);
    return await http.put(Uri.parse(url), headers: headers, body: jsonEncode(body));
  }
}
