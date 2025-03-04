package example.proto3;

@javax.annotation.Generated("by proto-qiu")
public final class Example {
public final static class StringInt32MapEntry extends com.protoc.qiu.GeneratedMessage {
    private java.lang.String key;
    private int value;

    public StringInt32MapEntry() {
    }

    public java.lang.String getKey() {
        return this.key;
    }

    public void setKey(java.lang.String key) {
        this.key = key;
    }

    public int getValue() {
        return this.value;
    }

    public void setValue(int value) {
        this.value = value;
    }

    public byte[] toByteArray() {
        java.io.ByteArrayOutputStream stream = new java.io.ByteArrayOutputStream();
        try {
            if (key != null) {
        writeString(stream, 1, key);
            }
            if (value != 0) {
        writeInt32(stream, 2, value);
            }
            return stream.toByteArray();
        } catch (Exception e) {
            throw new RuntimeException("Failed to serialize message", e);
        } finally {
            try {
                stream.close();
            } catch (Exception e) {
                // Ignore close exception
            }
        }
    }

    public static StringInt32MapEntry parseFrom(byte[] data) {
        java.io.ByteArrayInputStream stream = new java.io.ByteArrayInputStream(data);
        StringInt32MapEntry result = new StringInt32MapEntry();
        byte[] bytes;
        try {
            while (stream.available() > 0) {
                int tag = readTag(stream);
                int fieldNumber = getFieldNumberFromTag(tag);
                int wireType = getWireTypeFromTag(tag);
                switch (fieldNumber) {
                    case 1:
                    result.key = readString(stream);
                        break;
                    case 2:
                    result.value = readInt32(stream);
                        break;
                    default:
                        if (wireType == WIRETYPE_LENGTH_DELIMITED) {
                            readBytes(stream);
                        }
                        break;
                }
            }
        } catch (Exception e) {
            throw new RuntimeException("Failed to parse message", e);
        }
        return result;
    }
}
public final static class AllTypesDemo extends com.protoc.qiu.GeneratedMessage {
    private int int32Field;
    private long int64Field;
    private int uint32Field;
    private long uint64Field;
    private int sint32Field;
    private long sint64Field;
    private int fixed32Field;
    private long fixed64Field;
    private int sfixed32Field;
    private long sfixed64Field;
    private float floatField;
    private double doubleField;
    private boolean boolField;
    private java.lang.String stringField;
    private byte[] bytesField;
    private java.util.List<java.lang.Integer> repeatedInt32;
    private java.util.List<java.lang.String> repeatedString;
    private NestedMessage nestedMessage;
    private java.util.Map<java.lang.String, java.lang.Integer> mapField;
    private qiu.protobuf.Any anyField;
    private UserType userType;

    public AllTypesDemo() {
        this.repeatedInt32 = new java.util.ArrayList<>();
        this.repeatedString = new java.util.ArrayList<>();
        this.mapField = new java.util.HashMap<>();
    }

    public int getInt32Field() {
        return this.int32Field;
    }

    public void setInt32Field(int int32Field) {
        this.int32Field = int32Field;
    }

    public long getInt64Field() {
        return this.int64Field;
    }

    public void setInt64Field(long int64Field) {
        this.int64Field = int64Field;
    }

    public int getUint32Field() {
        return this.uint32Field;
    }

    public void setUint32Field(int uint32Field) {
        this.uint32Field = uint32Field;
    }

    public long getUint64Field() {
        return this.uint64Field;
    }

    public void setUint64Field(long uint64Field) {
        this.uint64Field = uint64Field;
    }

    public int getSint32Field() {
        return this.sint32Field;
    }

    public void setSint32Field(int sint32Field) {
        this.sint32Field = sint32Field;
    }

    public long getSint64Field() {
        return this.sint64Field;
    }

    public void setSint64Field(long sint64Field) {
        this.sint64Field = sint64Field;
    }

    public int getFixed32Field() {
        return this.fixed32Field;
    }

    public void setFixed32Field(int fixed32Field) {
        this.fixed32Field = fixed32Field;
    }

    public long getFixed64Field() {
        return this.fixed64Field;
    }

    public void setFixed64Field(long fixed64Field) {
        this.fixed64Field = fixed64Field;
    }

    public int getSfixed32Field() {
        return this.sfixed32Field;
    }

    public void setSfixed32Field(int sfixed32Field) {
        this.sfixed32Field = sfixed32Field;
    }

    public long getSfixed64Field() {
        return this.sfixed64Field;
    }

    public void setSfixed64Field(long sfixed64Field) {
        this.sfixed64Field = sfixed64Field;
    }

    public float getFloatField() {
        return this.floatField;
    }

    public void setFloatField(float floatField) {
        this.floatField = floatField;
    }

    public double getDoubleField() {
        return this.doubleField;
    }

    public void setDoubleField(double doubleField) {
        this.doubleField = doubleField;
    }

    public boolean getBoolField() {
        return this.boolField;
    }

    public void setBoolField(boolean boolField) {
        this.boolField = boolField;
    }

    public java.lang.String getStringField() {
        return this.stringField;
    }

    public void setStringField(java.lang.String stringField) {
        this.stringField = stringField;
    }

    public byte[] getBytesField() {
        return this.bytesField;
    }

    public void setBytesField(byte[] bytesField) {
        this.bytesField = bytesField;
    }

    public java.util.List<java.lang.Integer> getRepeatedInt32() {
        return this.repeatedInt32;
    }

    public void setRepeatedInt32(java.util.List<java.lang.Integer> repeatedInt32) {
        this.repeatedInt32 = repeatedInt32;
    }

    public java.util.List<java.lang.String> getRepeatedString() {
        return this.repeatedString;
    }

    public void setRepeatedString(java.util.List<java.lang.String> repeatedString) {
        this.repeatedString = repeatedString;
    }

    public NestedMessage getNestedMessage() {
        return this.nestedMessage;
    }

    public void setNestedMessage(NestedMessage nestedMessage) {
        this.nestedMessage = nestedMessage;
    }

    public java.util.Map<java.lang.String, java.lang.Integer> getMapField() {
        return this.mapField;
    }

    public void setMapField(java.util.Map<java.lang.String, java.lang.Integer> mapField) {
        this.mapField = mapField;
    }

    public qiu.protobuf.Any getAnyField() {
        return this.anyField;
    }

    public void setAnyField(qiu.protobuf.Any anyField) {
        this.anyField = anyField;
    }

    public UserType getUserType() {
        return this.userType;
    }

    public void setUserType(UserType userType) {
        this.userType = userType;
    }

    // OneOf: test_oneof
    private Object testOneof;
    private int testOneofCase = 0;
    public enum TestOneofCase {
        ONEOF_INT32(19),
        ONEOF_STRING(20),
        NOT_SET(0);
        private final int value;
        private TestOneofCase(int value) {
            this.value = value;
        }
        public int getValue() { return value; }
    }

    public int getOneofInt32() {
        if (testOneofCase == 19) {
            return (int) testOneof;
        }
        return 0;
    }

    public void setOneofInt32(int value) {
        testOneof = value;
        testOneofCase = 19;
    }

    public java.lang.String getOneofString() {
        if (testOneofCase == 20) {
            return (java.lang.String) testOneof;
        }
        return null;
    }

    public void setOneofString(java.lang.String value) {
        testOneof = value;
        testOneofCase = 20;
    }

    public TestOneofCase getTestOneofCase() {
        switch (testOneofCase) {
            case 19: return TestOneofCase.ONEOF_INT32;
            case 20: return TestOneofCase.ONEOF_STRING;
            default: return TestOneofCase.NOT_SET;
        }
    }

    public void clearTestOneof() {
        testOneof = null;
        testOneofCase = 0;
    }

public final static class NestedMessage extends com.protoc.qiu.GeneratedMessage {
    private int id;
    private java.lang.String name;

    public NestedMessage() {
    }

    public int getId() {
        return this.id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public java.lang.String getName() {
        return this.name;
    }

    public void setName(java.lang.String name) {
        this.name = name;
    }

    public byte[] toByteArray() {
        java.io.ByteArrayOutputStream stream = new java.io.ByteArrayOutputStream();
        try {
            if (id != 0) {
        writeInt32(stream, 1, id);
            }
            if (name != null) {
        writeString(stream, 2, name);
            }
            return stream.toByteArray();
        } catch (Exception e) {
            throw new RuntimeException("Failed to serialize message", e);
        } finally {
            try {
                stream.close();
            } catch (Exception e) {
                // Ignore close exception
            }
        }
    }

    public static NestedMessage parseFrom(byte[] data) {
        java.io.ByteArrayInputStream stream = new java.io.ByteArrayInputStream(data);
        NestedMessage result = new NestedMessage();
        byte[] bytes;
        try {
            while (stream.available() > 0) {
                int tag = readTag(stream);
                int fieldNumber = getFieldNumberFromTag(tag);
                int wireType = getWireTypeFromTag(tag);
                switch (fieldNumber) {
                    case 1:
                    result.id = readInt32(stream);
                        break;
                    case 2:
                    result.name = readString(stream);
                        break;
                    default:
                        if (wireType == WIRETYPE_LENGTH_DELIMITED) {
                            readBytes(stream);
                        }
                        break;
                }
            }
        } catch (Exception e) {
            throw new RuntimeException("Failed to parse message", e);
        }
        return result;
    }
}

    public byte[] toByteArray() {
        java.io.ByteArrayOutputStream stream = new java.io.ByteArrayOutputStream();
        try {
            if (int32Field != 0) {
        writeInt32(stream, 1, int32Field);
            }
            if (int64Field != 0L) {
        writeInt64(stream, 2, int64Field);
            }
            if (uint32Field != 0) {
        writeInt32(stream, 3, uint32Field);
            }
            if (uint64Field != 0L) {
        writeInt64(stream, 4, uint64Field);
            }
            if (sint32Field != 0) {
        writeSint32(stream, 5, sint32Field);
            }
            if (sint64Field != 0L) {
        writeSint64(stream, 6, sint64Field);
            }
            if (fixed32Field != 0) {
        writeFixed32(stream, 7, fixed32Field);
            }
            if (fixed64Field != 0L) {
        writeFixed64(stream, 8, fixed64Field);
            }
            if (sfixed32Field != 0) {
        writeFixed32(stream, 9, sfixed32Field);
            }
            if (sfixed64Field != 0L) {
        writeFixed64(stream, 10, sfixed64Field);
            }
            if (floatField != 0.0f) {
        writeFloat(stream, 11, floatField);
            }
            if (doubleField != 0.0) {
        writeDouble(stream, 12, doubleField);
            }
            if (boolField != false) {
        writeBool(stream, 13, boolField);
            }
            if (stringField != null) {
        writeString(stream, 14, stringField);
            }
            if (bytesField != null) {
        writeTag(stream, 15, WIRETYPE_LENGTH_DELIMITED);
        writeBytes(stream, bytesField);
            }
        if (repeatedInt32 != null) {
            for (java.lang.Integer item : repeatedInt32) {
        writeInt32(stream, 16, item);
            }
        }
        if (repeatedString != null) {
            for (java.lang.String item : repeatedString) {
        writeString(stream, 17, item);
            }
        }
            if (nestedMessage != null) {
            byte[] bytes = nestedMessage.toByteArray();
            writeTag(stream, 18, WIRETYPE_LENGTH_DELIMITED);
            writeBytes(stream, bytes);
            }
        if (mapField != null) {
            for (java.util.Map.Entry<java.lang.String, java.lang.Integer> entry : mapField.entrySet()) {
                writeTag(stream, 21, WIRETYPE_LENGTH_DELIMITED);
                StringInt32MapEntry e =new StringInt32MapEntry();
                e.setKey(entry.getKey());
                e.setValue(entry.getValue());
                byte[] mapBytes = e.toByteArray();
                writeBytes(stream, mapBytes);
            }
        }
            if (anyField != null) {
            byte[] bytes = anyField.toByteArray();
            writeTag(stream, 24, WIRETYPE_LENGTH_DELIMITED);
            writeBytes(stream, bytes);
            }
            if (userType != null) {
                writeInt32(stream, 23, userType.getNumber());
            }
            switch (testOneofCase) {
                case 19:
        writeInt32(stream, 19, (int)testOneof);
                    break;
                case 20:
        writeString(stream, 20, (java.lang.String)testOneof);
                    break;
            }
            return stream.toByteArray();
        } catch (Exception e) {
            throw new RuntimeException("Failed to serialize message", e);
        } finally {
            try {
                stream.close();
            } catch (Exception e) {
                // Ignore close exception
            }
        }
    }

    public static AllTypesDemo parseFrom(byte[] data) {
        java.io.ByteArrayInputStream stream = new java.io.ByteArrayInputStream(data);
        AllTypesDemo result = new AllTypesDemo();
        byte[] bytes;
        try {
            while (stream.available() > 0) {
                int tag = readTag(stream);
                int fieldNumber = getFieldNumberFromTag(tag);
                int wireType = getWireTypeFromTag(tag);
                switch (fieldNumber) {
                    case 1:
                    result.int32Field = readInt32(stream);
                        break;
                    case 2:
                    result.int64Field = readInt64(stream);
                        break;
                    case 3:
                    result.uint32Field = readInt32(stream);
                        break;
                    case 4:
                    result.uint64Field = readInt64(stream);
                        break;
                    case 5:
                    result.sint32Field = readSint32(stream);
                        break;
                    case 6:
                    result.sint64Field = readSint64(stream);
                        break;
                    case 7:
                    result.fixed32Field = readFixed32(stream);
                        break;
                    case 8:
                    result.fixed64Field = readFixed64(stream);
                        break;
                    case 9:
                    result.sfixed32Field = readFixed32(stream);
                        break;
                    case 10:
                    result.sfixed64Field = readFixed64(stream);
                        break;
                    case 11:
                    result.floatField = readFloat(stream);
                        break;
                    case 12:
                    result.doubleField = readDouble(stream);
                        break;
                    case 13:
                    result.boolField = readBool(stream);
                        break;
                    case 14:
                    result.stringField = readString(stream);
                        break;
                    case 15:
                    result.bytesField = readBytes(stream);
                        break;
                    case 16:
                    result.repeatedInt32.add(readInt32(stream));
                        break;
                    case 17:
                    result.repeatedString.add(readString(stream));
                        break;
                    case 18:
                    bytes = readBytes(stream);
                    result.nestedMessage = NestedMessage.parseFrom(bytes);
                        break;
                    case 21:
                    byte[] mapBytes = readBytes(stream);
                    java.io.ByteArrayInputStream mapStream = new java.io.ByteArrayInputStream(mapBytes);
                    java.lang.String key = null;
                    int value = 0;
                    while (mapStream.available() > 0) {
                        int mapTag = readTag(mapStream);
                        int mapFieldNumber = getFieldNumberFromTag(mapTag);
                        switch (mapFieldNumber) {
                            case 1: // key
                                key = readString(mapStream);
                                break;
                            case 2: // value
                                value = readInt32(mapStream);
                                break;
                            default:
                                break;
                        }
                    }
                    result.mapField.put(key, value);
                        break;
                    case 24:
                    bytes = readBytes(stream);
                    result.anyField = qiu.protobuf.Any.parseFrom(bytes);
                        break;
                    case 23:
                    result.userType = UserType.forNumber(readInt32(stream));
                        break;
                    case 19: // oneof test_oneof
                        result.testOneof = readInt32(stream);
                        result.testOneofCase = 19;
                        break;
                    case 20: // oneof test_oneof
                        result.testOneof = readString(stream);
                        result.testOneofCase = 20;
                        break;
                    default:
                        if (wireType == WIRETYPE_LENGTH_DELIMITED) {
                            readBytes(stream);
                        }
                        break;
                }
            }
        } catch (Exception e) {
            throw new RuntimeException("Failed to parse message", e);
        }
        return result;
    }
}
public enum UserType {
    UNKNOWN(0),
    ADMIN(1),
    GUEST(2);

    private final int value;

    UserType(int value) {
        this.value = value;
    }

    public int getNumber() {
        return value;
    }
 	public static UserType forNumber(int value) {
        for (UserType e : values()) {
            if (e.value == value) {
                return e;
            }
        }
        return null;
    }
}
}
