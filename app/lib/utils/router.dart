import 'package:go_router/go_router.dart';
import '../models/vacancy_model.dart';
import '../models/vacancy_resume_match_model.dart';
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
    path: RouteNames.vacancy,
    builder: (context, state) {
      final vacancy = state.extra as VacancyModel;
      return VacancyScreen(vacancy: vacancy);
    },
  ),
  GoRoute(
    path: RouteNames.vacancyResumeMatch,
    builder: (context, state) {
      final resumeMatch = state.extra as VacancyResumeMatchModel;
      return VacancyResumeMatchScreen(resumeMatch: resumeMatch);
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
