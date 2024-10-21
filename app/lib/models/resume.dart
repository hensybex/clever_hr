// models/resume.dart

import 'package:hive/hive.dart';
import 'hive_type_ids.dart';

part 'resume.g.dart';

@HiveType(typeId: HiveTypeIds.resume)
class Resume extends HiveObject {
  @HiveField(0)
  int? id;

  @HiveField(1)
  String uploadedFrom;

  @HiveField(2)
  String pdfPath;

  @HiveField(3)
  String rawText;

  @HiveField(4)
  String cleanText;

  @HiveField(5)
  String standarizedText;

  @HiveField(6)
  int? jobGroupId;

  @HiveField(7)
  int? specializationId;

  @HiveField(8)
  int? qualificationId;

  @HiveField(9)
  String? fullName;

  @HiveField(10)
  int? gptCallId;

  @HiveField(11)
  DateTime? createdAt;

  @HiveField(12)
  DateTime? updatedAt;

  Resume({
    this.id,
    required this.uploadedFrom,
    required this.pdfPath,
    required this.rawText,
    required this.cleanText,
    required this.standarizedText,
    this.jobGroupId,
    this.specializationId,
    this.qualificationId,
    this.fullName,
    this.gptCallId,
    this.createdAt,
    this.updatedAt,
  });

  // Factory method to create an instance from JSON
  factory Resume.fromJson(Map<String, dynamic> json) {
    print("------------PRINTING KEYS------------");
    print(json.keys); // Prints top-level keys
    print(json['resume']?.keys.join(', ')); // Prints all keys inside 'resume' without '...'
    print("------------PRINTING KEYS------------");

    final resumeData = json['resume']; // Access the nested 'resume' object

    return Resume(
      id: resumeData['ID'],
      uploadedFrom: resumeData['UploadedFrom'] ?? '',
      pdfPath: resumeData['PDFPath'] ?? '',
      rawText: resumeData['RawText'] ?? '',
      cleanText: resumeData['CleanText'] ?? '',
      standarizedText: resumeData['StandarizedText'] ?? '',
      jobGroupId: resumeData['JobGroupID'],
      specializationId: resumeData['SpecializationID'],
      qualificationId: resumeData['QualificationID'],
      fullName: resumeData['FullName'] as String?,
      gptCallId: resumeData['GPTCallID'],
      createdAt: resumeData['CreatedAt'] != null ? DateTime.parse(resumeData['CreatedAt']) : null,
      updatedAt: resumeData['UpdatedAt'] != null ? DateTime.parse(resumeData['UpdatedAt']) : null,
    );
  }

  // Method to convert the instance to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'uploadedFrom': uploadedFrom,
      'pdfPath': pdfPath,
      'rawText': rawText,
      'cleanText': cleanText,
      'standarizedText': standarizedText,
      'jobGroupId': jobGroupId,
      'specializationId': specializationId,
      'qualificationId': qualificationId,
      'fullName': fullName,
      'gptCallId': gptCallId,
      'createdAt': createdAt?.toIso8601String(),
      'updatedAt': updatedAt?.toIso8601String(),
    };
  }
}
