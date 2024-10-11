// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'vacancy_resume_match_model.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class VacancyResumeMatchModelAdapter
    extends TypeAdapter<VacancyResumeMatchModel> {
  @override
  final int typeId = 1;

  @override
  VacancyResumeMatchModel read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return VacancyResumeMatchModel(
      id: fields[0] as String,
      vacancyId: fields[1] as String,
      resumeId: fields[2] as String,
      suitability: fields[3] as double,
      analysedForMatch: fields[4] as double,
    );
  }

  @override
  void write(BinaryWriter writer, VacancyResumeMatchModel obj) {
    writer
      ..writeByte(5)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.vacancyId)
      ..writeByte(2)
      ..write(obj.resumeId)
      ..writeByte(3)
      ..write(obj.suitability)
      ..writeByte(4)
      ..write(obj.analysedForMatch);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is VacancyResumeMatchModelAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
