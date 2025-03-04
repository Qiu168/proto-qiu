package qiu.protobuf;

import com.protoc.qiu.GeneratedMessage;

public final class Any extends GeneratedMessage {
    private String typeUrl;
    private byte[] value;

    public Any() {
        this.typeUrl = "";
        this.value = new byte[0];
    }

    public String getTypeUrl() {
        return typeUrl;
    }

    public void setTypeUrl(String typeUrl) {
        this.typeUrl = typeUrl;
    }

    public byte[] getValue() {
        return value;
    }

    public void setValue(byte[] value) {
        this.value = value;
    }

    public static Any pack(GeneratedMessage message) {
        Any any = new Any();
        any.typeUrl = message.getClass().getName();
        any.value = message.toByteArray();
        return any;
    }

    public <T extends GeneratedMessage> T unpack(Class<T> clazz) {
        try {
            String expectedType = clazz.getName();
            if (!typeUrl.equals(expectedType)) {
                throw new RuntimeException("Type mismatch: expected " + expectedType + " but got " + typeUrl);
            }
            return (T)clazz.getMethod("parseFrom", byte[].class).invoke(null, value);
        } catch (Exception e) {
            throw new RuntimeException("Failed to unpack Any message", e);
        }
    }

    public boolean is(Class<? extends GeneratedMessage> clazz) {
        String expectedType = clazz.getName();
        return typeUrl.equals(expectedType);
    }

    public byte[] toByteArray() {
        java.io.ByteArrayOutputStream stream = new java.io.ByteArrayOutputStream();
        try {
            // Write typeUrl
            writeString(stream, 1, typeUrl);
            // Write value
            writeTag(stream, 2, WIRETYPE_LENGTH_DELIMITED);
            writeBytes(stream, value);
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

    public static Any parseFrom(byte[] data) {
        java.io.ByteArrayInputStream stream = new java.io.ByteArrayInputStream(data);
        Any result = new Any();
        try {
            while (stream.available() > 0) {
                int tag = readTag(stream);
                int fieldNumber = getFieldNumberFromTag(tag);
                switch (fieldNumber) {
                    case 1:
                        result.typeUrl = readString(stream);
                        break;
                    case 2:
                        result.value = readBytes(stream);
                        break;
                    default:
                        if (getWireTypeFromTag(tag) == WIRETYPE_LENGTH_DELIMITED) {
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