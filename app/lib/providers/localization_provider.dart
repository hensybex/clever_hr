// providers/localization_provider.dart

import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';

class LocalizationProvider extends ChangeNotifier {
  final FlutterLocalization _localization = FlutterLocalization.instance;

  String currentLanguage = "Russian";

  void switchLanguage(String languageCode) {
    _localization.translate(languageCode);
    currentLanguage = languageCode;
    notifyListeners();
  }

  String get currentLanguageName {
    switch (currentLanguage) {
      case 'en':
        return 'English';
      case 'ru':
        return 'Russian';
      default:
        return 'English';
    }
  }
}
