// models/user_token.dart

import 'package:hive/hive.dart';

import 'hive_type_ids.dart';

part 'user_token.g.dart';

@HiveType(typeId: HiveTypeIds.userToken)
class UserToken extends HiveObject {
  @HiveField(0)
  String token;

  @HiveField(1)
  DateTime expiryDate;

  UserToken({required this.token, required this.expiryDate});
}
