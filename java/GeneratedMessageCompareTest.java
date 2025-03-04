package com.protoc.qiu;

import com.google.protobuf.CodedOutputStream;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.*;

public class GeneratedMessageCompareTest {

    @Test
    public void compareVarint32Encoding() throws Exception {
        int[] testValues = {0, 1, -1, 127, 128, 16383, 16384, Integer.MAX_VALUE, Integer.MIN_VALUE};

        for (int value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeVarint32(ourStream, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeInt32NoTag(value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

//            System.out.println(Arrays.toString(ourBytes));
//            System.out.println(Arrays.toString(protobufBytes));

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Varint32 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareVarint64Encoding() throws Exception {
        long[] testValues = {0L, 1L, -1L, 127L, 128L, 16383L, 16384L, Long.MAX_VALUE, Long.MIN_VALUE};

        for (long value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeVarint64(ourStream, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeUInt64NoTag(value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Varint64 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareStringEncoding() throws Exception {
        String[] testValues = {"", "hello", "世界", "Hello World!", "🌟"};

        for (String value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeString(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeString(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("String encoding mismatch for value '%s'", value));
        }
    }

    @Test
    public void compareFixed32Encoding() throws Exception {
        int[] testValues = {0, 1, -1, Integer.MAX_VALUE, Integer.MIN_VALUE};

        for (int value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeFixed32(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeFixed32(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Fixed32 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareFixed64Encoding() throws Exception {
        long[] testValues = {0L, 1L, -1L, Long.MAX_VALUE, Long.MIN_VALUE};

        for (long value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeFixed64(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeFixed64(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Fixed64 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareSInt32Encoding() throws Exception {
        int[] testValues = {0, 1, -1, Integer.MAX_VALUE, Integer.MIN_VALUE};

        for (int value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeSint32(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeSInt32(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("SInt32 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareSInt64Encoding() throws Exception {
        long[] testValues = {0L, 1L, -1L, Long.MAX_VALUE, Long.MIN_VALUE};

        for (long value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeSint64(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream
                    .newInstance(byteArrayOutputStream);
            codedOutputStream.writeSInt64(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("SInt64 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareFloatEncoding() throws Exception {
        float[] testValues = {0.0f, 1.0f, -1.0f, Float.MAX_VALUE, Float.MIN_VALUE, Float.NaN, Float.POSITIVE_INFINITY, Float.NEGATIVE_INFINITY};

        for (float value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeFloat(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeFloat(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Float encoding mismatch for value %f", value));
        }
    }

    @Test
    public void compareDoubleEncoding() throws Exception {
        double[] testValues = {0.0, 1.0, -1.0, Double.MAX_VALUE, Double.MIN_VALUE, Double.NaN, Double.POSITIVE_INFINITY, Double.NEGATIVE_INFINITY};

        for (double value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeDouble(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeDouble(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Double encoding mismatch for value %f", value));
        }
    }

    @Test
    public void compareBoolEncoding() throws Exception {
        boolean[] testValues = {true, false};

        for (boolean value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeBool(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeBool(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Bool encoding mismatch for value %b", value));
        }
    }

    @Test
    public void compareBytesEncoding() throws Exception {
        byte[][] testValues = {
                new byte[]{},
                new byte[]{1, 2, 3},
                new byte[]{-1, -2, -3},
                new byte[]{Byte.MAX_VALUE, Byte.MIN_VALUE}
        };

        for (byte[] value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeBytes(ourStream, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeByteArrayNoTag(value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Bytes encoding mismatch for value %s", Arrays.toString(value)));
        }
    }

    @Test
    public void compareEnumEncoding() throws Exception {
        int[] testValues = {0, 1, 2, 100, -1};

        for (int value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeEnum(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeEnum(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("Enum encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareSFixed32Encoding() throws Exception {
        int[] testValues = {0, 1, -1, Integer.MAX_VALUE, Integer.MIN_VALUE};

        for (int value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeSFixed32(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeSFixed32(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("SFixed32 encoding mismatch for value %d", value));
        }
    }

    @Test
    public void compareSFixed64Encoding() throws Exception {
        long[] testValues = {0L, 1L, -1L, Long.MAX_VALUE, Long.MIN_VALUE};

        for (long value : testValues) {
            // 我们的实现
            ByteArrayOutputStream ourStream = new ByteArrayOutputStream();
            GeneratedMessage.writeSFixed64(ourStream, 1, value);
            byte[] ourBytes = ourStream.toByteArray();

            // Protobuf的实现
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            CodedOutputStream codedOutputStream = CodedOutputStream.newInstance(byteArrayOutputStream);
            codedOutputStream.writeSFixed64(1, value);
            codedOutputStream.flush();
            byte[] protobufBytes = byteArrayOutputStream.toByteArray();

            assertArrayEquals(protobufBytes, ourBytes,
                    String.format("SFixed64 encoding mismatch for value %d", value));
        }
    }
}