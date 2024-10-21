// lib/widgets/expandable_text.dart

import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import '../utils/locales.dart';
class ExpandableTextWidget extends StatefulWidget {
  final String label;
  final String text;
  const ExpandableTextWidget({
    super.key,
    required this.label,
    required this.text,
  });
  @override
  ExpandableTextWidgetState createState() => ExpandableTextWidgetState();
}
class ExpandableTextWidgetState extends State<ExpandableTextWidget> {
  bool _isExpanded = false;
  @override
  Widget build(BuildContext context) {
    final firstRow = widget.text.split('\n').first;
    return InkWell(
      onTap: () => setState(() => _isExpanded = !_isExpanded), // Toggle expand/collapse
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '${widget.label}:',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
            Text(
              _isExpanded ? widget.text : firstRow,
              maxLines: _isExpanded ? null : 1,
              overflow: _isExpanded ? TextOverflow.visible : TextOverflow.ellipsis,
            ),
            Text(
              _isExpanded ? AppLocale.showLess.getString(context) : AppLocale.showMore.getString(context),
              style: const TextStyle(color: Colors.blue),
            ),
          ],
        ),
      ),
    );
  }
}
