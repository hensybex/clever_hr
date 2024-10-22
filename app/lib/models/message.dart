class Message {
  final DateTime createdAt;
  final String tgLogin;
  final String message;
  final bool sentByUser;

  Message({
    required this.createdAt,
    required this.tgLogin,
    required this.message,
    required this.sentByUser,
  });

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      createdAt: DateTime.parse(json['created_at']),
      tgLogin: json['tg_login'],
      message: json['message'],
      sentByUser: json['sent_by_user'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'created_at': createdAt.toIso8601String(),
      'tg_login': tgLogin,
      'message': message,
      'sent_by_user': sentByUser,
    };
  }
}
