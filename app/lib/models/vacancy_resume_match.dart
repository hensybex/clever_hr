// lib/models/vacancy_resume_match.dart

import 'package:hive/hive.dart';
import 'hive_type_ids.dart';
part 'vacancy_resume_match.g.dart';
@HiveType(typeId: HiveTypeIds.vacancyResumeMatch)
class VacancyResumeMatch extends HiveObject {
  @HiveField(0)
  int id;
  @HiveField(1)
  int vacancyId;
  @HiveField(2)
  int resumeId;
  @HiveField(3)
  double suitability; // Value between 0.0 and 1.0 (Matches `Score` in Go model)
  @HiveField(4)
  double analysedForMatch; // Value between 0.0 and 1.0
  @HiveField(5)
  AnalysisField relevantWorkExperience;
  @HiveField(6)
  AnalysisField technicalSkillsAndProficiencies;
  @HiveField(7)
  AnalysisField educationAndCertifications;
  @HiveField(8)
  AnalysisField softSkillsAndCulturalFit;
  @HiveField(9)
  AnalysisField languageAndCommunicationSkills;
  @HiveField(10)
  AnalysisField problemSolvingAndAnalyticalAbilities;
  @HiveField(11)
  AnalysisField adaptabilityAndLearningCapacity;
  @HiveField(12)
  AnalysisField leadershipAndManagementExperience;
  @HiveField(13)
  AnalysisField motivationAndCareerObjectives;
  @HiveField(14)
  AnalysisField additionalQualificationsAndValueAdds;
  @HiveField(15)
  String analysisStatus;
  @HiveField(16)
  DateTime createdAt;
  VacancyResumeMatch({
    required this.id,
    required this.vacancyId,
    required this.resumeId,
    required this.suitability,
    required this.analysedForMatch,
    required this.relevantWorkExperience,
    required this.technicalSkillsAndProficiencies,
    required this.educationAndCertifications,
    required this.softSkillsAndCulturalFit,
    required this.languageAndCommunicationSkills,
    required this.problemSolvingAndAnalyticalAbilities,
    required this.adaptabilityAndLearningCapacity,
    required this.leadershipAndManagementExperience,
    required this.motivationAndCareerObjectives,
    required this.additionalQualificationsAndValueAdds,
    required this.analysisStatus,
    required this.createdAt,
  });
  // Factory method to create an instance from JSON
  factory VacancyResumeMatch.fromJson(Map<String, dynamic> json) {
    return VacancyResumeMatch(
      id: json['ID'], // Convert to String
      vacancyId: json['VacancyID'], // Convert to String
      resumeId: json['ResumeID'], // Convert to String
      suitability: json['Score'].toDouble(),
      analysedForMatch: json['analysedForMatch']?.toDouble() ?? 0.0,
      relevantWorkExperience: AnalysisField.fromJson(json['RelevantWorkExperience']),
      technicalSkillsAndProficiencies: AnalysisField.fromJson(json['TechnicalSkillsAndProficiencies']),
      educationAndCertifications: AnalysisField.fromJson(json['EducationAndCertifications']),
      softSkillsAndCulturalFit: AnalysisField.fromJson(json['SoftSkillsAndCulturalFit']),
      languageAndCommunicationSkills: AnalysisField.fromJson(json['LanguageAndCommunicationSkills']),
      problemSolvingAndAnalyticalAbilities: AnalysisField.fromJson(json['ProblemSolvingAndAnalyticalAbilities']),
      adaptabilityAndLearningCapacity: AnalysisField.fromJson(json['AdaptabilityAndLearningCapacity']),
      leadershipAndManagementExperience: AnalysisField.fromJson(json['LeadershipAndManagementExperience']),
      motivationAndCareerObjectives: AnalysisField.fromJson(json['MotivationAndCareerObjectives']),
      additionalQualificationsAndValueAdds: AnalysisField.fromJson(json['AdditionalQualificationsAndValueAdds']),
      analysisStatus: json['AnalysisStatus'] ?? "", // Handle nullable status
      createdAt: DateTime.parse(json['CreatedAt']),
    );
  }
  // Method to convert the instance to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'vacancyId': vacancyId,
      'resumeId': resumeId,
      'suitability': suitability,
      'analysedForMatch': analysedForMatch,
      'relevantWorkExperience': relevantWorkExperience.toJson(),
      'technicalSkillsAndProficiencies': technicalSkillsAndProficiencies.toJson(),
      'educationAndCertifications': educationAndCertifications.toJson(),
      'softSkillsAndCulturalFit': softSkillsAndCulturalFit.toJson(),
      'languageAndCommunicationSkills': languageAndCommunicationSkills.toJson(),
      'problemSolvingAndAnalyticalAbilities': problemSolvingAndAnalyticalAbilities.toJson(),
      'adaptabilityAndLearningCapacity': adaptabilityAndLearningCapacity.toJson(),
      'leadershipAndManagementExperience': leadershipAndManagementExperience.toJson(),
      'motivationAndCareerObjectives': motivationAndCareerObjectives.toJson(),
      'additionalQualificationsAndValueAdds': additionalQualificationsAndValueAdds.toJson(),
      'analysisStatus': analysisStatus,
      'createdAt': createdAt.toIso8601String(),
    };
  }
}
@HiveType(typeId: HiveTypeIds.analysisField)
class AnalysisField {
  @HiveField(0)
  String overview;
  @HiveField(1)
  int score;
  AnalysisField({
    required this.overview,
    required this.score,
  });
  // Factory method to create an AnalysisField instance from JSON
  factory AnalysisField.fromJson(Map<String, dynamic> json) {
    return AnalysisField(
      overview: json['overview'],
      score: json['score'],
    );
  }
  // Method to convert the AnalysisField to JSON
  Map<String, dynamic> toJson() {
    return {
      'overview': overview,
      'score': score,
    };
  }
}
