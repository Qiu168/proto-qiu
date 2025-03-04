package example.proto3;

import org.junit.jupiter.api.Test;
import qiu.protobuf.Any;

import static org.junit.jupiter.api.Assertions.*;

public class ExampleTest {
    
    @Test
    public void testAllTypesDemo() throws Exception {
        // 创建测试对象
        Example.AllTypesDemo original = new Example.AllTypesDemo();
        
        // 设置基本类型字段
        original.setInt32Field(-123);
        original.setInt64Field(-12345678900L);
        original.setUint32Field(123);
        original.setUint64Field(12345678900L);
        original.setSint32Field(-456);
        original.setSint64Field(-45678900L);
        original.setFixed32Field(789);
        original.setFixed64Field(7890123456L);
        original.setSfixed32Field(-789);
        original.setSfixed64Field(-7890123456L);
        original.setFloatField(3.14f);
        original.setDoubleField(3.14159265359);
        original.setBoolField(true);
        original.setStringField("Hello, Proto!");
        original.setBytesField(new byte[]{1, 2, 3, 4, 5});

        // 设置重复字段
        original.getRepeatedInt32().add(1);
        original.getRepeatedInt32().add(2);
        original.getRepeatedInt32().add(3);
        original.getRepeatedString().add("one");
        original.getRepeatedString().add("two");
        original.getRepeatedString().add("three");

        // 设置嵌套消息
        Example.AllTypesDemo.NestedMessage nested = new Example.AllTypesDemo.NestedMessage();
        nested.setId(100);
        nested.setName("Nested Message");
        original.setNestedMessage(nested);

        // 设置 Map 字段
        original.getMapField().put("key1", 1);
        original.getMapField().put("key2", 2);

        // 设置枚举字段
        original.setUserType(Example.UserType.ADMIN);

        // 设置 oneof 字段
        original.setOneofString("oneof test");

        original.setAnyField(Any.pack(nested));

        // 序列化
        byte[] serialized = original.toByteArray();

        // 反序列化
        Example.AllTypesDemo parsed = Example.AllTypesDemo.parseFrom(serialized);

        // 验证所有字段
        assertEquals(original.getInt32Field(), parsed.getInt32Field());
        assertEquals(original.getInt64Field(), parsed.getInt64Field());
        assertEquals(original.getUint32Field(), parsed.getUint32Field());
        assertEquals(original.getUint64Field(), parsed.getUint64Field());
        assertEquals(original.getSint32Field(), parsed.getSint32Field());
        assertEquals(original.getSint64Field(), parsed.getSint64Field());
        assertEquals(original.getFixed32Field(), parsed.getFixed32Field());
        assertEquals(original.getFixed64Field(), parsed.getFixed64Field());
        assertEquals(original.getSfixed32Field(), parsed.getSfixed32Field());
        assertEquals(original.getSfixed64Field(), parsed.getSfixed64Field());
        assertEquals(original.getFloatField(), parsed.getFloatField(), 0.0001);
        assertEquals(original.getDoubleField(), parsed.getDoubleField(), 0.0001);
        assertEquals(original.getBoolField(), parsed.getBoolField());
        assertEquals(original.getStringField(), parsed.getStringField());
        assertArrayEquals(original.getBytesField(), parsed.getBytesField());
        
        // 验证重复字段
        assertEquals(original.getRepeatedInt32(), parsed.getRepeatedInt32());
        assertEquals(original.getRepeatedString(), parsed.getRepeatedString());
        
        // 验证嵌套消息
        assertEquals(original.getNestedMessage().getId(), parsed.getNestedMessage().getId());
        assertEquals(original.getNestedMessage().getName(), parsed.getNestedMessage().getName());
        
        // 验证 Map 字段
        assertEquals(original.getMapField(), parsed.getMapField());
        
        // 验证枚举字段
        assertEquals(original.getUserType(), parsed.getUserType());
        
        // 验证 oneof 字段
        assertEquals(original.getOneofString(), parsed.getOneofString());
        assertEquals(Example.AllTypesDemo.TestOneofCase.ONEOF_STRING, parsed.getTestOneofCase());
    }

    @Test
    public void testMapEntry() throws Exception {
        Example.StringInt32MapEntry original = new Example.StringInt32MapEntry();
        original.setKey("test_key");
        original.setValue(42);

        byte[] serialized = original.toByteArray();
        Example.StringInt32MapEntry parsed = Example.StringInt32MapEntry.parseFrom(serialized);

        assertEquals(original.getKey(), parsed.getKey());
        assertEquals(original.getValue(), parsed.getValue());
    }

    @Test
    public void testEnum() {
        assertEquals(0, Example.UserType.UNKNOWN.getNumber());
        assertEquals(1, Example.UserType.ADMIN.getNumber());
        assertEquals(2, Example.UserType.GUEST.getNumber());

        assertEquals(Example.UserType.UNKNOWN, Example.UserType.forNumber(0));
        assertEquals(Example.UserType.ADMIN, Example.UserType.forNumber(1));
        assertEquals(Example.UserType.GUEST, Example.UserType.forNumber(2));
        assertNull(Example.UserType.forNumber(99));
    }
}