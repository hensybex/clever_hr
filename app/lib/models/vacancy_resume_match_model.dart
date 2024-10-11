import 'package:hive/hive.dart';

part 'vacancy_resume_match_model.g.dart';

@HiveType(typeId: 1)
class VacancyResumeMatchModel extends HiveObject {
  @HiveField(0)
  String id;

  @HiveField(1)
  String vacancyId;

  @HiveField(2)
  String resumeId;

  @HiveField(3)
  double suitability; // Value between 0.0 and 1.0

  @HiveField(4)
  double analysedForMatch; // Value between 0.0 and 1.0

  VacancyResumeMatchModel({
    required this.id,
    required this.vacancyId,
    required this.resumeId,
    required this.suitability,
    required this.analysedForMatch,
  });
}
