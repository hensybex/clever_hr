// models/vacancy.dart

import 'package:hive/hive.dart';
import 'hive_type_ids.dart';

part 'vacancy.g.dart';

@HiveType(typeId: HiveTypeIds.vacancy)
class Vacancy extends HiveObject {
  @HiveField(0)
  int? id;

  @HiveField(1)
  String description;

  @HiveField(2)
  String status;

  @HiveField(3)
  String? title;

  @HiveField(4)
  String? standarizedText;

  @HiveField(5)
  int? jobGroupId;

  Vacancy({
    this.id,
    required this.description,
    this.status = 'BeingAnalysed',
    this.title,
    this.standarizedText,
    this.jobGroupId,
  });

  // Factory method to create an instance from JSON
  factory Vacancy.fromJson(Map<String, dynamic> json) {
    return Vacancy(
      id: json['ID'], // Handle nullable id
      description: json['Description'] ?? '', // Fallback to empty string if description is null
      status: json['status'] ?? 'BeingAnalysed', // Handle nullable status with default value
      title: json['Title'] as String?, // Handle nullable title
      standarizedText: json['standarizedText'] as String?, // Handle nullable standarizedText
      jobGroupId: json['JobGroupID'], // Handle nullable jobGroupId
    );
  }

  // Method to convert the instance to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'description': description,
      'status': status,
      'title': title,
      'standarizedText': standarizedText,
      'jobGroupId': jobGroupId,
    };
  }
}
