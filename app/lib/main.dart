import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_localization/flutter_localization.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'services/hive/hive_init.dart';
import 'utils/locales.dart';
import 'utils/router.dart';
import 'providers/auth_provider.dart';
import 'providers/vacancy_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await dotenv.load(fileName: "assets/.env");
  await HiveInitializer.instance.initHive();
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> with AppLocale {
  final FlutterLocalization localization = FlutterLocalization.instance;

  @override
  void initState() {
    super.initState();

    String systemLanguageCode = PlatformDispatcher.instance.locale.languageCode;

    localization.init(
      mapLocales: [
        const MapLocale('en', AppLocale.EN),
        const MapLocale('ru', AppLocale.RU),
      ],
      initLanguageCode: systemLanguageCode,
    );

    localization.onTranslatedLanguage = _onTranslatedLanguage;
  }

  void _onTranslatedLanguage(Locale? locale) {
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProvider(create: (_) => VacancyProvider()),
      ],
      child: const App(),
    );
  }
}

class App extends StatefulWidget {
  const App({super.key});

  @override
  AppState createState() => AppState();
}

class AppState extends State<App> {
  late final GoRouter router;

  @override
  void initState() {
    super.initState();
    router = GoRouter(
      routes: routes,
    );
  }

  @override
  Widget build(BuildContext context) {
    final localization = FlutterLocalization.instance;

    return MaterialApp.router(
      routerConfig: router,
      title: 'HR App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      darkTheme: ThemeData.dark(),
      themeMode: ThemeMode.system,
      supportedLocales: localization.supportedLocales,
      localizationsDelegates: localization.localizationsDelegates,
      locale: localization.currentLocale,
      localeResolutionCallback: (locale, supportedLocales) {
        for (var supportedLocale in supportedLocales) {
          if (supportedLocale.languageCode == locale?.languageCode) {
            return supportedLocale;
          }
        }
        return const Locale('en');
      },
    );
  }
}