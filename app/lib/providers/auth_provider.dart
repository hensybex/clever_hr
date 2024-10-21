// lib/providers/auth_provider.dart

import 'package:flutter/material.dart';
import 'package:hive/hive.dart';
import '../services/api/auth_service.dart';
import '../models/user_token.dart';
import '../services/hive/hive_init.dart';
class AuthProvider with ChangeNotifier {
  final AuthService _authService;
  AuthProvider(this._authService);
  bool _isAuthenticated = false;
  bool get isAuthenticated => _isAuthenticated;
  Future<void> login(String username, String password) async {
    try {
      // Call the login method from AuthService
      var response = await _authService.login(username, password);
      // Get the token and expiry date from the response
      String token = response['token'];
      DateTime expiryDate = DateTime.parse(response['expire']);
      // Save the token and expiry date in Hive
      var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
      tokenBox.put('token', UserToken(token: token, expiryDate: expiryDate));
      _isAuthenticated = true;
      notifyListeners();
    } catch (e) {
      // Handle the error appropriately
      print("Login failed: ${e.toString()}");
      _isAuthenticated = false;
      notifyListeners();
    }
  }
  void logout() async {
    try {
      // Call the logout method from AuthService
      await _authService.logout();
      // Clear the token from Hive
      var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
      tokenBox.clear();
      _isAuthenticated = false;
      notifyListeners();
    } catch (e) {
      print("Logout failed: ${e.toString()}");
    }
  }
  String? getToken() {
    var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
    UserToken? userToken = tokenBox.get('token');
    return userToken?.token;
  }
}
