// lib/services/hive_init.dart
import 'package:app/models/vacancy_model.dart';
import 'package:app/models/vacancy_resume_match_model.dart';
import 'package:hive_flutter/hive_flutter.dart';

import '../../utils/config.dart';

class HiveInitializer {
  static final HiveInitializer _instance = HiveInitializer._privateConstructor();
  HiveInitializer._privateConstructor();
  static HiveInitializer get instance => _instance;

  ClearHiveBoxes clearHiveBoxes = ClearHiveBoxes.doNotClear;

  Future<void> initHive() async {
    await Hive.initFlutter();
    Hive.registerAdapter(VacancyModelAdapter());
    Hive.registerAdapter(VacancyResumeMatchModelAdapter());

    if (clearHiveBoxes == ClearHiveBoxes.clear) {
      await clearAllBoxes();
    }

    await Hive.openBox<VacancyModel>(BoxNames.vacancies);
    await Hive.openBox<VacancyResumeMatchModel>(BoxNames.vacancyResumesMatches);
  }

  Future<void> clearAllBoxes() async {
    await Hive.deleteBoxFromDisk(BoxNames.vacancies);
    await Hive.deleteBoxFromDisk(BoxNames.vacancyResumesMatches);
  }
}

class BoxNames {
  static String vacancies = 'vacancies';
  static String vacancyResumesMatches = 'vacancyResumesMatches';
}
