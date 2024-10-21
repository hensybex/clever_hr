// providers/vacancy_provider.dart

import 'dart:async';

import 'package:app/services/api/websocket_service.dart';
import 'package:app/services/hive/hive_init.dart';
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import '../models/server_message.dart';
import '../models/vacancy.dart';
import '../models/vacancy_resume_match.dart';
import '../services/api/vacancy_service.dart';
import 'package:hive/hive.dart';
import '../services/api/match_service.dart';

class VacancyProvider with ChangeNotifier {
  final VacancyService _vacancyService;
  final MatchService _matchService;
  final WebSocketService _wsService;
  VacancyProvider(this._vacancyService, this._matchService, this._wsService);
  List<Vacancy> _vacancies = [];
  Vacancy? _currentVacancy;
  StreamSubscription<ServerMessage>? _wsSubscription;
  List<VacancyResumeMatch> _vacancyResumeMatches = [];

  List<Vacancy> get vacancies => _vacancies;
  Vacancy? get vacancy => _currentVacancy;
  List<VacancyResumeMatch> get vacancyResumeMatches => _vacancyResumeMatches;

  bool isLoading = false;
  String errorMessage = '';
  String waitingMessage = '';
  String displayedContent = '';

  Future<void> loadVacancies() async {
    try {
      var box = Hive.box<Vacancy>(BoxNames.vacancies);
      _vacancies = box.values.toList();

      var fetchedVacancies = await _vacancyService.getVacancies();

      _vacancies = fetchedVacancies;

      await box.clear();
      await box.addAll(_vacancies);

      notifyListeners(); // Notify that the data has been updated
    } catch (e) {
      print("Error while loading vacancies: $e");
    }
  }

  Future<void> fetchVacancyResumeMatches(int vacancyId) async {
    isLoading = true;

    // Defer notifyListeners to prevent state changes during build
    SchedulerBinding.instance.addPostFrameCallback((_) {
      notifyListeners();
    });

    try {
      final fetchedResumes = await _matchService.fetchVacancyResumeMatches(vacancyId);
      _vacancyResumeMatches = fetchedResumes;
      isLoading = false;
    } catch (e) {
      errorMessage = 'Error fetching resume matches: $e';
      isLoading = false;
    }

    // Defer notifyListeners to prevent state changes during build
    SchedulerBinding.instance.addPostFrameCallback((_) {
      notifyListeners();
    });
  }

  Future<void> fetchVacancy(int vacancyId) async {
    try {
      final Vacancy fetchedVacancy = await _vacancyService.getVacancyByID(vacancyId);
      _currentVacancy = fetchedVacancy;
    } catch (e) {
      errorMessage = 'Error fetching resume matches: $e';
      isLoading = false;
    }
    notifyListeners();
  }

  Future<void> createVacancy(Vacancy vacancy) async {
    try {
      // Log the beginning of the WebSocket connection attempt
      print("Starting WebSocket connection for vacancy creation...");

      isLoading = true;
      errorMessage = '';
      waitingMessage = 'Uploading vacancy...';
      notifyListeners();

      // Start the WebSocket connection
      Stream<ServerMessage> messageStream = _wsService.createVacancyWebSocket(vacancy);

      // Log successful WebSocket connection initialization
      print("WebSocket connection initialized successfully");

      // Listen to the WebSocket stream
      _wsSubscription?.cancel();
      _wsSubscription = messageStream.listen(
        (serverMessage) async {
          // Log every received server message
          print("Received server message: ${serverMessage.status}");

          if (serverMessage.status != null) {
            // Handle status updates
            waitingMessage = serverMessage.status!;
            print("Status update: $waitingMessage");
            notifyListeners();
          }

          if (serverMessage.result != null) {
            // Handle result content
            displayedContent += serverMessage.result!;
            print("Received result chunk: ${serverMessage.result!}");
            notifyListeners();
          }

          if (serverMessage.status == "Vacancy processing and matching completed") {
            // Update vacancy status
            vacancy.status = 'Analysed';
            print("Vacancy processing completed, updating status.");

            // Add the vacancy to the list and Hive cache
            _vacancies.add(vacancy);
            var box = Hive.box<Vacancy>(BoxNames.vacancies);
            await box.add(vacancy);
            print("Vacancy added to Hive cache successfully");

            isLoading = false;
            waitingMessage = '';
            notifyListeners();
          }
        },
        onDone: () {
          print("WebSocket connection closed.");
          isLoading = false;
          waitingMessage = '';
          notifyListeners();
        },
        onError: (error) {
          // Log WebSocket errors
          print("WebSocket error: $error");
          errorMessage = 'Error: $error';
          isLoading = false;
          waitingMessage = '';
          notifyListeners();
        },
      );
    } catch (e) {
      // Log any exceptions thrown during WebSocket setup or connection
      print("Exception caught during WebSocket operation: ${e.toString()}");
      errorMessage = 'Error: ${e.toString()}';
      isLoading = false;
      waitingMessage = '';
      notifyListeners();
    }
  }

  @override
  void dispose() {
    _wsSubscription?.cancel();
    super.dispose();
  }
}
