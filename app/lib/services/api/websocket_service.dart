// services/api/websocket_service.dart

import 'dart:convert';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:web_socket_channel/html.dart';

import '../../models/server_message.dart';
import '../../models/user_token.dart';
import '../../models/vacancy.dart';
import '../hive/hive_init.dart';
import 'api_client.dart';

class WebSocketService {
  final ApiClient apiClient;

  WebSocketService(this.apiClient);

  Stream<ServerMessage> createVacancyWebSocket(Vacancy vacancy) async* {
    var tokenBox = Hive.box<UserToken>(BoxNames.userToken);
    UserToken? userToken = tokenBox.get('token');
    if (userToken == null || userToken.token.isEmpty) {
      throw Exception("No valid token found.");
    }

    final wsScheme = apiClient.baseUrl.startsWith('https') ? 'wss' : 'ws';
    final wsUri = '$wsScheme://${Uri.parse(apiClient.baseUrl).host}:8080/api/vacancies/upload?token=${userToken.token}';

    try {
      final channel = HtmlWebSocketChannel.connect(wsUri);

      final Map<String, dynamic> message = vacancy.toJson();
      channel.sink.add(json.encode(message));

      yield* channel.stream.map((event) {
        final data = json.decode(event);
        return ServerMessage.fromJson(data);
      }).handleError((error) {});
    } catch (e) {
      rethrow;
    }
  }
}
