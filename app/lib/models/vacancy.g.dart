// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'vacancy.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class VacancyAdapter extends TypeAdapter<Vacancy> {
  @override
  final int typeId = 1;

  @override
  Vacancy read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Vacancy(
      id: fields[0] as int?,
      description: fields[1] as String,
      status: fields[2] as String,
      title: fields[3] as String?,
      standarizedText: fields[4] as String?,
      jobGroupId: fields[5] as int?,
    );
  }

  @override
  void write(BinaryWriter writer, Vacancy obj) {
    writer
      ..writeByte(6)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.description)
      ..writeByte(2)
      ..write(obj.status)
      ..writeByte(3)
      ..write(obj.title)
      ..writeByte(4)
      ..write(obj.standarizedText)
      ..writeByte(5)
      ..write(obj.jobGroupId);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is VacancyAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
