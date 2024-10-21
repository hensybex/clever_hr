// services/hive/hive_init.dart

import 'package:app/models/resume.dart';
import 'package:app/models/user_token.dart';
import 'package:hive_flutter/hive_flutter.dart';

import '../../models/vacancy.dart';
import '../../models/vacancy_resume_match.dart';
import '../../utils/config.dart';

class HiveInitializer {
  static final HiveInitializer _instance = HiveInitializer._privateConstructor();
  HiveInitializer._privateConstructor();
  static HiveInitializer get instance => _instance;

  ClearHiveBoxes clearHiveBoxes = ClearHiveBoxes.doNotClear;

  Future<void> initHive() async {
    await Hive.initFlutter();
    Hive.registerAdapter(VacancyAdapter());
    Hive.registerAdapter(ResumeAdapter());
    Hive.registerAdapter(VacancyResumeMatchAdapter());
    Hive.registerAdapter(UserTokenAdapter());

    if (clearHiveBoxes == ClearHiveBoxes.clear) {
      await clearAllBoxes();
    }

    await Hive.openBox<Vacancy>(BoxNames.vacancies);
    await Hive.openBox<Resume>(BoxNames.resumes);
    await Hive.openBox<VacancyResumeMatch>(BoxNames.vacancyResumesMatches);
    await Hive.openBox<UserToken>(BoxNames.userToken);
  }

  Future<void> clearAllBoxes() async {
    await Hive.deleteBoxFromDisk(BoxNames.vacancies);
    await Hive.deleteBoxFromDisk(BoxNames.resumes);
    await Hive.deleteBoxFromDisk(BoxNames.vacancyResumesMatches);
    await Hive.deleteBoxFromDisk(BoxNames.userToken);
  }
}

class BoxNames {
  static String vacancies = 'vacancies';
  static String resumes = 'resumes';
  static String vacancyResumesMatches = 'vacancyResumesMatches';
  static String userToken = 'userToken';
}
