// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'vacancy_resume_match.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class VacancyResumeMatchAdapter extends TypeAdapter<VacancyResumeMatch> {
  @override
  final int typeId = 2;

  @override
  VacancyResumeMatch read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return VacancyResumeMatch(
      id: fields[0] as int,
      vacancyId: fields[1] as int,
      resumeId: fields[2] as int,
      suitability: fields[3] as double,
      analysedForMatch: fields[4] as double,
      relevantWorkExperience: fields[5] as AnalysisField,
      technicalSkillsAndProficiencies: fields[6] as AnalysisField,
      educationAndCertifications: fields[7] as AnalysisField,
      softSkillsAndCulturalFit: fields[8] as AnalysisField,
      languageAndCommunicationSkills: fields[9] as AnalysisField,
      problemSolvingAndAnalyticalAbilities: fields[10] as AnalysisField,
      adaptabilityAndLearningCapacity: fields[11] as AnalysisField,
      leadershipAndManagementExperience: fields[12] as AnalysisField,
      motivationAndCareerObjectives: fields[13] as AnalysisField,
      additionalQualificationsAndValueAdds: fields[14] as AnalysisField,
      analysisStatus: fields[15] as String,
      createdAt: fields[16] as DateTime,
    );
  }

  @override
  void write(BinaryWriter writer, VacancyResumeMatch obj) {
    writer
      ..writeByte(17)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.vacancyId)
      ..writeByte(2)
      ..write(obj.resumeId)
      ..writeByte(3)
      ..write(obj.suitability)
      ..writeByte(4)
      ..write(obj.analysedForMatch)
      ..writeByte(5)
      ..write(obj.relevantWorkExperience)
      ..writeByte(6)
      ..write(obj.technicalSkillsAndProficiencies)
      ..writeByte(7)
      ..write(obj.educationAndCertifications)
      ..writeByte(8)
      ..write(obj.softSkillsAndCulturalFit)
      ..writeByte(9)
      ..write(obj.languageAndCommunicationSkills)
      ..writeByte(10)
      ..write(obj.problemSolvingAndAnalyticalAbilities)
      ..writeByte(11)
      ..write(obj.adaptabilityAndLearningCapacity)
      ..writeByte(12)
      ..write(obj.leadershipAndManagementExperience)
      ..writeByte(13)
      ..write(obj.motivationAndCareerObjectives)
      ..writeByte(14)
      ..write(obj.additionalQualificationsAndValueAdds)
      ..writeByte(15)
      ..write(obj.analysisStatus)
      ..writeByte(16)
      ..write(obj.createdAt);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) || other is VacancyResumeMatchAdapter && runtimeType == other.runtimeType && typeId == other.typeId;
}

class AnalysisFieldAdapter extends TypeAdapter<AnalysisField> {
  @override
  final int typeId = 3;

  @override
  AnalysisField read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return AnalysisField(
      overview: fields[0] as String,
      score: fields[1] as int,
    );
  }

  @override
  void write(BinaryWriter writer, AnalysisField obj) {
    writer
      ..writeByte(2)
      ..writeByte(0)
      ..write(obj.overview)
      ..writeByte(1)
      ..write(obj.score);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) || other is AnalysisFieldAdapter && runtimeType == other.runtimeType && typeId == other.typeId;
}
