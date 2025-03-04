package com.protoc.qiu;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;

public abstract class GeneratedMessage {
    public static final int WIRETYPE_VARINT = 0;
    public static final int WIRETYPE_FIXED64 = 1;
    public static final int WIRETYPE_LENGTH_DELIMITED = 2;
    public static final int WIRETYPE_START_GROUP = 3;
    public static final int WIRETYPE_END_GROUP = 4;
    public static final int WIRETYPE_FIXED32 = 5;
  
    public static void writeTag(ByteArrayOutputStream stream, int fieldNumber, int wireType){
        writeVarint32(stream, (fieldNumber << 3) | wireType);
    }
    public static void writeString(ByteArrayOutputStream stream,int fieldNumber, String str) throws IOException {
        writeTag(stream, fieldNumber, WIRETYPE_LENGTH_DELIMITED);
        writeBytes(stream, str.getBytes());
    }
    public static void writeVarint32(ByteArrayOutputStream stream, int value){
        while (true) {
            if ((value & ~0x7F) == 0) {
                stream.write(value);
                return;
            } else {
                stream.write((value & 0x7F) | 0x80);
                value >>>= 7;
            }
        }
    }   
    public static void writeInt32(ByteArrayOutputStream stream, int fieldNumber, int value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint32(stream, value);
    }
    public static void writeInt64(ByteArrayOutputStream stream, int fieldNumber, long value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint64(stream, value);
    }
    public static void writeUint32(ByteArrayOutputStream stream, int fieldNumber, int value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint32(stream, value);
    }
    
    public static void writeUint64(ByteArrayOutputStream stream, int fieldNumber, long value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint64(stream, value);
    }
    public static void writeBytes(ByteArrayOutputStream stream, byte[] bytes) throws IOException {
        writeVarint32(stream, bytes.length);
        stream.write(bytes);
    }
    public static void writeVarint64(ByteArrayOutputStream stream, long value){
        while (true) {
            if ((value & ~0x7F) == 0) {
                stream.write((int) value);
                return;
            } else {
                stream.write(((int) value & 0x7F) | 0x80);
                value >>>= 7;
            }
        }
    }
    public static void writeFloat(ByteArrayOutputStream stream, int fieldNumber, float value){
        writeFixed32(stream,fieldNumber, Float.floatToIntBits(value));
    }
    public static void writeDouble(ByteArrayOutputStream stream, int fieldNumber, double value){
        writeFixed64(stream,fieldNumber, Double.doubleToLongBits(value));
    }
    public static void writeBool(ByteArrayOutputStream stream, int fieldNumber, boolean value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        stream.write(value ? 1 : 0);
    }
    
    public static void writeEnum(ByteArrayOutputStream stream, int fieldNumber, int value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint32(stream, value);
    }
    public static void writeSint32(ByteArrayOutputStream stream, int fieldNumber, int value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint32(stream, (value << 1) ^ (value >> 31));
    }
    
    public static void writeSint64(ByteArrayOutputStream stream, int fieldNumber, long value){
        writeTag(stream, fieldNumber, WIRETYPE_VARINT);
        writeVarint64(stream, (value << 1) ^ (value >> 63));
    }
    public static void writeFixed32(ByteArrayOutputStream stream, int fieldNumber, int value){
        writeTag(stream, fieldNumber, WIRETYPE_FIXED32);
        stream.write(value >> 24);
        stream.write(value >> 16);
        stream.write(value >> 8);
        stream.write(value);
    }
    
    public static void writeFixed64(ByteArrayOutputStream stream, int fieldNumber, long value){
        writeTag(stream, fieldNumber, WIRETYPE_FIXED64);
        stream.write((int) (value >> 56));
        stream.write((int) (value >> 48));
        stream.write((int) (value >> 40));
        stream.write((int) (value >> 32));
        stream.write((int) (value >> 24));
        stream.write((int) (value >> 16));
        stream.write((int) (value >> 8));
        stream.write((int) value);
    }
    public static void writeSFixed32(ByteArrayOutputStream stream, int fieldNumber, int value){
        writeFixed32(stream,fieldNumber, value);
    }
    public static void writeSFixed64(ByteArrayOutputStream stream, int fieldNumber, long value){
        writeFixed64(stream,fieldNumber, value);
    }

    public static int readVarint32(ByteArrayInputStream stream){
        int value = 0;
        int shift = 0;
        while (true) {
            int b = stream.read();
            if (b == -1) {
                throw new RuntimeException("Malformed varint");
            }
            value |= (b & 0x7F) << shift;
            if ((b & 0x80) == 0) {
                return value;
            }
            shift += 7;
            if (shift >= 32) {
                throw new RuntimeException("Malformed varint");
            }
        }
    }

    public static long readVarint64(ByteArrayInputStream stream){
        long value = 0;
        int shift = 0;
        while (true) {
            int b = stream.read();
            if (b == -1) {
                throw new RuntimeException("Malformed varint");
            }
            value |= (long)(b & 0x7F) << shift;
            if ((b & 0x80) == 0) {
                return value;
            }
            shift += 7;
            if (shift >= 64) {
                throw new RuntimeException("Malformed varint");
            }
        }
    }

    public static byte[] readBytes(ByteArrayInputStream stream) throws IOException {
        int length = readVarint32(stream);
        byte[] bytes = new byte[length];
        if (stream.read(bytes)!= length) {
            throw new RuntimeException("Malformed bytes");
        }
        return bytes;
    }

    public static int readTag(ByteArrayInputStream stream){
        int tag = readVarint32(stream);
        return tag;
    }
    
    public static int getFieldNumberFromTag(int tag) {
        return tag >>> 3;
    }
    
    public static int getWireTypeFromTag(int tag) {
        return tag & 0x7;
    }
    public static String readString(ByteArrayInputStream stream) throws IOException {
        int length = readVarint32(stream);
        byte[] bytes = new byte[length];
        if (stream.read(bytes) != length) {
            throw new RuntimeException("Malformed string");
        }
        return new String(bytes);
    }
    public static int readInt32(ByteArrayInputStream stream){
        return readVarint32(stream);
    }
    public static long readInt64(ByteArrayInputStream stream){
        return readVarint64(stream);
    }
    public static float readFloat(ByteArrayInputStream stream) throws IOException {
        return Float.intBitsToFloat(readFixed32(stream));
    }
    public static double readDouble(ByteArrayInputStream stream) throws IOException {
        return Double.longBitsToDouble(readFixed64(stream));
    }
    public static boolean readBool(ByteArrayInputStream stream){
        return readVarint32(stream) != 0;
    }
    public static int readEnum(ByteArrayInputStream stream){
        return readInt32(stream);
    }
    public static int readSint32(ByteArrayInputStream stream){
        int n = readInt32(stream);
        return (n >>> 1) ^ -(n & 1);
    }
    public static long readSint64(ByteArrayInputStream stream){
        long n = readInt64(stream);
        return (n >>> 1) ^ -(n & 1);
    }
    public static int readFixed32(ByteArrayInputStream stream) throws IOException {
        byte[] bytes = new byte[4];
        if (stream.read(bytes) != 4) {
            throw new RuntimeException("Malformed fixed32");
        }
        return ((bytes[0] & 0xFF) << 24) 
             | ((bytes[1] & 0xFF) << 16) 
             | ((bytes[2] & 0xFF) << 8) 
             | (bytes[3] & 0xFF);
    }
    
    public static long readFixed64(ByteArrayInputStream stream) throws IOException {
        byte[] bytes = new byte[8];
        if (stream.read(bytes) != 8) {
            throw new RuntimeException("Malformed fixed64");
        }
        return ((long)(bytes[0] & 0xFF) << 56)
             | ((long)(bytes[1] & 0xFF) << 48)
             | ((long)(bytes[2] & 0xFF) << 40)
             | ((long)(bytes[3] & 0xFF) << 32)
             | ((long)(bytes[4] & 0xFF) << 24)
             | ((long)(bytes[5] & 0xFF) << 16)
             | ((long)(bytes[6] & 0xFF) << 8)
             | (bytes[7] & 0xFF);
    }
    public static int readSFixed32(ByteArrayInputStream stream) throws IOException {
        return readFixed32(stream);
    }
    public static long readSFixed64(ByteArrayInputStream stream) throws IOException {
        return readFixed64(stream);
    }
    public static int readUint32(ByteArrayInputStream stream){
        return readInt32(stream);
    }
    public static long readUint64(ByteArrayInputStream stream){
        return readInt64(stream);
    }

    public abstract byte[] toByteArray();
    
    // static
    // public abstract GeneratedMessage parseFrom(byte[] bytes);
}