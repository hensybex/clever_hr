import 'package:hive/hive.dart';

part 'vacancy_model.g.dart';

@HiveType(typeId: 0)
class VacancyModel extends HiveObject {
  @HiveField(0)
  String id;

  @HiveField(1)
  String description;

  @HiveField(2)
  String priceRange;

  @HiveField(3)
  int yearsOfExperience;

  @HiveField(4)
  String workType; // 'Remote' or 'Office'

  @HiveField(5)
  String status; // 'BeingAnalysed' or 'Analysed'

  VacancyModel({
    required this.id,
    required this.description,
    required this.priceRange,
    required this.yearsOfExperience,
    required this.workType,
    this.status = 'BeingAnalysed',
  });
}
