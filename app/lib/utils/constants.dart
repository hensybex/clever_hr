// lib/utils/constants.dart

import 'package:flutter_dotenv/flutter_dotenv.dart';

final String apiBaseUrl = dotenv.env['API_BASE_URL'] ?? 'http://localhost:8080/';

// Log the API base URL for debugging purposes
void logApiBaseUrl() {
  print('API Base URL: $apiBaseUrl');
}
