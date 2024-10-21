// providers/match_provider.dart

import 'dart:async';

import 'package:app/services/api/resume_service.dart';
import 'package:app/services/hive/hive_init.dart';
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import '../models/resume.dart';
import '../models/server_message.dart';
import '../models/vacancy.dart';
import '../models/vacancy_resume_match.dart';
import '../services/api/vacancy_service.dart';
import 'package:hive/hive.dart';
import '../services/api/match_service.dart';

class MatchProvider with ChangeNotifier {
  final VacancyService _vacancyService;
  final ResumeService _resumeService;
  final MatchService _matchService;
  MatchProvider(this._vacancyService, this._resumeService, this._matchService);
  List<VacancyResumeMatch> _vacancyResumeMatches = [];
  Vacancy? _vacancy;
  Resume? _resume;
  VacancyResumeMatch? _match;

  Vacancy? get vacancy => _vacancy;
  Resume? get resume => _resume;
  VacancyResumeMatch? get match => _match;
  List<VacancyResumeMatch> get vacancyResumeMatches => _vacancyResumeMatches;

  Future<void> loadResume(int resumeID) async {
    try {
      var fetchedResume = await _resumeService.getResumeByID(resumeID);
      _resume = fetchedResume;

      notifyListeners();
    } catch (e) {
      print("Error while loading vacancies: $e");
    }
  }

  Future<VacancyResumeMatch?> loadVacancyResumeMatch(int matchId) async {
    try {
      // Fetch the VacancyResumeMatch using matchId
      var currentResumeMatch = await _matchService.fetchVacancyResumeMatchById(matchId);
      _match = currentResumeMatch;
      notifyListeners();
    } catch (e) {
      print("Error while loading match: $e");
    }
    return _match;
  }

  void updateMatch(VacancyResumeMatch VRMatch) {
    _match = VRMatch;
    notifyListeners();
  }
}
