// lib/models/server_message.dart

class ServerMessage {
  final String? status;
  final String? result;
  final String? details;
  ServerMessage({this.status, this.result, this.details});
  factory ServerMessage.fromJson(Map<String, dynamic> json) {
    return ServerMessage(
      status: json['status'],
      result: json['result'],
      details: json['details'],
    );
  }
}
