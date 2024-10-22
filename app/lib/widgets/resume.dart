// lib/widgets/resume.dart

import 'package:flutter/material.dart';
import 'package:flutter_localization/flutter_localization.dart';
import '../models/resume.dart';
import '../services/api/api_client.dart';
import '../services/danya_api_service.dart';
import '../utils/constants.dart';
import '../utils/locales.dart';
import 'expandable_text.dart';
class ResumeWidget extends StatelessWidget {
  final Resume resume;
  final int matchID;
  const ResumeWidget({super.key, required this.resume, required this.matchID});
  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          AppLocale.name.getString(context),
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        Row(
          children: [
            Text('${resume.fullName}'),
            ElevatedButton(
              onPressed: () async {
                // Create the ApiClient and DanyaApiService
                final apiClient = ApiClient(apiBaseUrl);
                final danyaApiService = DanyaApiService(apiClient);

                try {
                  // Call the service to generate and send the message
                  await danyaApiService.generateAndSendMessage(matchID, '@hensybex2');
                  print('Message sent successfully');
                } catch (e) {
                  print('Error: $e');
                }
              },
              child: const Text("Press me"),
            ),
          ],
        ),
        const SizedBox(height: 10),
        // Displaying uploadedFrom field
        Text(
          AppLocale.uploadedFrom.getString(context),
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 10),
        // Reusing ExpandableTextWidget for Clean Text
        ExpandableTextWidget(
          label: AppLocale.cleanText.getString(context),
          text: resume.cleanText,
        ),
        const SizedBox(height: 10),
        // Reusing ExpandableTextWidget for Standardized Text
        ExpandableTextWidget(
          label: AppLocale.standardizedText.getString(context),
          text: resume.standarizedText,
        ),
        const SizedBox(height: 10),
        // Upload resume button
        ElevatedButton(
          onPressed: () {
            // Handle resume upload logic
            print('Upload resume button clicked');
          },
          child: Text(AppLocale.uploadResume.getString(context)),
        ),
      ],
    );
  }
}
