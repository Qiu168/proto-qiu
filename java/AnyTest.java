package qiu.protobuf;

import com.protoc.qiu.GeneratedMessage;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class AnyTest {

    @Test
    public void testPackAndUnpack() {

        TestMessage original = new TestMessage();
        original.setName("test");
        original.setId(123);

        Any any = Any.pack(original);
        Assertions.assertNotNull(any);
        assertEquals(TestMessage.class.getName(), any.getTypeUrl());
        assertNotNull(any.getValue());
        assertTrue(any.getValue().length > 0);

        TestMessage unpacked = any.unpack(TestMessage.class);
        assertNotNull(unpacked);
        assertEquals(original.getName(), unpacked.getName());
        assertEquals(original.getId(), unpacked.getId());
    }

    @Test
    public void testIs() {
        TestMessage msg = new TestMessage();
        Any any = Any.pack(msg);

        assertTrue(any.is(TestMessage.class));
        assertFalse(any.is(AnotherTestMessage.class));
    }

    @Test
    public void testUnpackTypeMismatch() {
        TestMessage msg = new TestMessage();
        Any any = Any.pack(msg);
        assertThrowsExactly(RuntimeException.class, () -> any.unpack(AnotherTestMessage.class));
    }

    @Test
    public void testSerializeAndParse() {
        TestMessage original = new TestMessage();
        original.setName("test");
        original.setId(123);

        Any any = Any.pack(original);
        byte[] serialized = any.toByteArray();
        assertNotNull(serialized);
        assertTrue(serialized.length > 0);

        Any parsed = Any.parseFrom(serialized);
        assertNotNull(parsed);
        assertEquals(any.getTypeUrl(), parsed.getTypeUrl());
        assertArrayEquals(any.getValue(), parsed.getValue());

        TestMessage unpacked = parsed.unpack(TestMessage.class);
        assertEquals(original.getName(), unpacked.getName());
        assertEquals(original.getId(), unpacked.getId());
    }

    // 用于测试的示例消息类
    public static class TestMessage extends GeneratedMessage {
        private String name = "";
        private int id = 0;

        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }

        public int getId() {
            return id;
        }

        public void setId(int id) {
            this.id = id;
        }

        public byte[] toByteArray() {
            java.io.ByteArrayOutputStream stream = new java.io.ByteArrayOutputStream();
            try {
                writeString(stream, 1, name);
                writeInt32(stream, 2, id);
                return stream.toByteArray();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        }

        public static TestMessage parseFrom(byte[] data) {
            java.io.ByteArrayInputStream stream = new java.io.ByteArrayInputStream(data);
            TestMessage result = new TestMessage();
            try {
                while (stream.available() > 0) {
                    int tag = readTag(stream);
                    int fieldNumber = getFieldNumberFromTag(tag);
                    switch (fieldNumber) {
                        case 1:
                            result.name = readString(stream);
                            break;
                        case 2:
                            result.id = readInt32(stream);
                            break;
                        default:
                            if (getWireTypeFromTag(tag) == WIRETYPE_LENGTH_DELIMITED) {
                                readBytes(stream);
                            }
                            break;
                    }
                }
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
            return result;
        }
    }

    // 用于类型不匹配测试的另一个消息类
    public static class AnotherTestMessage extends com.protoc.qiu.GeneratedMessage {
        @Override
        public byte[] toByteArray() {
            return new byte[0];
        }

        public static AnotherTestMessage parseFrom(byte[] data) {
            return new AnotherTestMessage();
        }
    }
}