// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'resume.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class ResumeAdapter extends TypeAdapter<Resume> {
  @override
  final int typeId = 4;

  @override
  Resume read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Resume(
      id: fields[0] as int?,
      uploadedFrom: fields[1] as String,
      pdfPath: fields[2] as String,
      rawText: fields[3] as String,
      cleanText: fields[4] as String,
      standarizedText: fields[5] as String,
      jobGroupId: fields[6] as int?,
      specializationId: fields[7] as int?,
      qualificationId: fields[8] as int?,
      fullName: fields[9] as String?,
      gptCallId: fields[10] as int?,
      createdAt: fields[11] as DateTime?,
      updatedAt: fields[12] as DateTime?,
    );
  }

  @override
  void write(BinaryWriter writer, Resume obj) {
    writer
      ..writeByte(13)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.uploadedFrom)
      ..writeByte(2)
      ..write(obj.pdfPath)
      ..writeByte(3)
      ..write(obj.rawText)
      ..writeByte(4)
      ..write(obj.cleanText)
      ..writeByte(5)
      ..write(obj.standarizedText)
      ..writeByte(6)
      ..write(obj.jobGroupId)
      ..writeByte(7)
      ..write(obj.specializationId)
      ..writeByte(8)
      ..write(obj.qualificationId)
      ..writeByte(9)
      ..write(obj.fullName)
      ..writeByte(10)
      ..write(obj.gptCallId)
      ..writeByte(11)
      ..write(obj.createdAt)
      ..writeByte(12)
      ..write(obj.updatedAt);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) || other is ResumeAdapter && runtimeType == other.runtimeType && typeId == other.typeId;
}
