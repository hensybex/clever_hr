// lib/utils/router.dart

import 'package:go_router/go_router.dart';
import '../screens/login_screen.dart';
import '../screens/home_screen.dart';
import '../screens/create_vacancy_screen.dart';
import '../screens/vacancy_screen.dart';
import '../screens/vacancy_resume_match_screen.dart';
final routes = [
  GoRoute(
    path: RouteNames.login,
    builder: (context, state) => LoginScreen(),
  ),
  GoRoute(
    path: RouteNames.home,
    builder: (context, state) => const HomeScreen(),
  ),
  GoRoute(
    path: RouteNames.createVacancy,
    builder: (context, state) => const CreateVacancyScreen(),
  ),
  GoRoute(
    path: '${RouteNames.vacancy}/:id',
    builder: (context, state) {
      final id = int.parse(state.pathParameters['id']!);
      return VacancyScreen(vacancyId: id);
    },
  ),
  GoRoute(
    path: '${RouteNames.vacancyResumeMatch}/:matchId',
    builder: (context, state) {
      final matchId = int.parse(state.pathParameters['matchId']!);
      return VacancyResumeMatchScreen(matchId: matchId);
    },
  ),
];
class RouteNames {
  static String login = '/';
  static String home = '/home';
  static String createVacancy = '/create-vacancy';
  static String vacancy = '/vacancy';
  static String vacancyResumeMatch = '/vacancy-resume-match';
}
