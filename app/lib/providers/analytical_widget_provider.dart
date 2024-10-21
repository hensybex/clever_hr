import 'package:flutter/material.dart';

class AnalyticalWidgetProvider with ChangeNotifier {
  // List of hover states for 10 AnalyticalWidgets
  List<bool> _hoverStates = List.generate(10, (_) => false);

  List<bool> get hoverStates => _hoverStates;

  // Toggle hover state for a specific index
  void setHoverState(int index, bool isHovered) {
    _hoverStates[index] = isHovered;
    notifyListeners();
  }

  bool getHoverState(int index) {
    return _hoverStates[index];
  }
}
