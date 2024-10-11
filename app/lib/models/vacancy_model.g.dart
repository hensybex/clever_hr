// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'vacancy_model.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class VacancyModelAdapter extends TypeAdapter<VacancyModel> {
  @override
  final int typeId = 0;

  @override
  VacancyModel read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return VacancyModel(
      id: fields[0] as String,
      description: fields[1] as String,
      priceRange: fields[2] as String,
      yearsOfExperience: fields[3] as int,
      workType: fields[4] as String,
      status: fields[5] as String,
    );
  }

  @override
  void write(BinaryWriter writer, VacancyModel obj) {
    writer
      ..writeByte(6)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.description)
      ..writeByte(2)
      ..write(obj.priceRange)
      ..writeByte(3)
      ..write(obj.yearsOfExperience)
      ..writeByte(4)
      ..write(obj.workType)
      ..writeByte(5)
      ..write(obj.status);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is VacancyModelAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
