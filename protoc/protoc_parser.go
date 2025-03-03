// 语法分析器
package protoc

import (
	"fmt"
	"io"
	"proto-qiu/constant"
	"strconv"
)

type Parser struct {
	lexer        *Lexer
	protoc       *Protoc
	currentToken Token
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		lexer:  NewLexer(r),
		protoc: &Protoc{},
	}
}

func (p *Parser) expectAndAdvance(t TokenType, values ...string) error {
	if err := p.expect(t, values...); err != nil {
		return err
	}
	return p.advance()
}

// expect 函数用于检查当前 token 是否符合预期的类型和值
func (p *Parser) expect(t TokenType, values ...string) error {
	if p.currentToken.Type != t {
		return fmt.Errorf("expected %v, got %v", t, p.currentToken.Type)
	}
	if len(values) > 0 {
		for _, v := range values {
			if p.currentToken.Value == v {
				return nil
			}
		}
		return fmt.Errorf("expected one of %v, got %v", values, p.currentToken.Value)
	}
	return nil
}

// advance 函数用于前进到下一个 token
func (p *Parser) advance() error {
	token, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.currentToken = token
	return nil
}

func (p *Parser) Parse() (*Protoc, error) {
	if err := p.advance(); err != nil {
		return nil, err
	}

	for p.currentToken.Type != TokenEOF {
		switch p.currentToken.Value {
		case constant.KeywordSyntax:
			if err := p.parseSyntax(); err != nil {
				return nil, err
			}
		case constant.KeywordPackage:
			if err := p.parsePackage(); err != nil {
				return nil, err
			}
		case constant.KeywordImport:
			if err := p.parseImport(); err != nil {
				return nil, err
			}
		case constant.KeywordMessage:
			msg, err := p.parseMessage()
			if err != nil {
				return nil, err
			}
			p.protoc.Messages = append(p.protoc.Messages, msg)
		case constant.KeywordEnum:
			enum, err := p.parseEnum()
			if err != nil {
				return nil, err
			}
			p.protoc.Enums = append(p.protoc.Enums, enum)
		case constant.KeywordService:
			service, err := p.parseService()
			if err != nil {
				return nil, err
			}
			p.protoc.Services = append(p.protoc.Services, service)
		default:
			return nil, fmt.Errorf("unexpected token: %v", p.currentToken.Value)
		}
		if p.currentToken.Value == ";" {
			if err := p.advance(); err != nil {
			}
		}
	}

	return p.protoc, nil
}

func (p *Parser) parseSyntax() error {
	if err := p.advance(); err != nil { // 跳过 'syntax'
		return err
	}
	if err := p.expect(TokenSymbol, constant.SymbolEqual); err != nil {
		return err
	}
	if err := p.advance(); err != nil { // 跳过 '='
		return err
	}
	if err := p.expect(TokenString); err != nil {
		return err
	}
	p.protoc.SyntaxVersion = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过字符串
		return err
	}
	return p.expect(TokenSymbol, constant.SymbolSemicolon)
}

func (p *Parser) parsePackage() error {
	if err := p.advance(); err != nil { // 跳过 'package'
		return err
	}

	for p.currentToken.Value != constant.SymbolSemicolon {
		p.protoc.PackageName += p.currentToken.Value
		if err := p.advance(); err != nil {
			return err
		}
	}

	return p.expect(TokenSymbol, constant.SymbolSemicolon)
}

func (p *Parser) parseImport() error {
	imp := &Import{}
	if err := p.advance(); err != nil { // 跳过 'import'
		return err
	}
	if p.currentToken.Value == constant.KeywordPublic {
		imp.Public = true
		if err := p.advance(); err != nil {
			return err
		}
	}
	if err := p.expect(TokenString); err != nil {
		return err
	}
	imp.Path = p.currentToken.Value
	p.protoc.Imports = append(p.protoc.Imports, imp)
	if err := p.advance(); err != nil { // 跳过字符串
		return err
	}
	return p.expect(TokenSymbol, constant.SymbolSemicolon)
}

func (p *Parser) parseMessage() (*Message, error) {
	msg := &Message{}
	// 跳过 'message'
	if err := p.advance(); err != nil {
		return nil, err
	}
	if err := p.expect(TokenIdent); err != nil {
		return nil, err
	}
	msg.Name = p.currentToken.Value
	// 跳过消息名
	if err := p.advance(); err != nil {
		return nil, err
	}
	if err := p.expect(TokenSymbol, constant.SymbolLeftBrace); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '{'
		return nil, err
	}

	// 解析消息体内容
	for p.currentToken.Value != constant.SymbolRightBrace {
		switch p.currentToken.Value {
		case constant.KeywordMessage:
			nestedMsg, err := p.parseMessage()
			if err != nil {
				return nil, err
			}
			msg.InnerMessages = append(msg.InnerMessages, nestedMsg)
		case constant.KeywordEnum:
			enum, err := p.parseEnum()
			if err != nil {
				return nil, err
			}
			enum.SuperMessage = msg
			msg.Enums = append(msg.Enums, enum)
		case constant.KeywordOneof:
			oneof, err := p.parseOneOf()
			if err != nil {
				return nil, err
			}
			msg.OneOfs = append(msg.OneOfs, oneof)
		default:
			field, err := p.parseField()
			if err != nil {
				return nil, err
			}
			msg.Fields = append(msg.Fields, field)
		}
	}

	if err := p.advance(); err != nil { // 跳过 '}'
		return nil, err
	}
	return msg, nil
}

func (p *Parser) parseField() (*Field, error) {
	field := &Field{Options: &FieldOptions{}}
	if p.currentToken.Value == constant.KeywordMap {
		if err := p.advance(); err != nil {
			return nil, err
		}
		err := p.expectAndAdvance(TokenSymbol, constant.SymbolLessThan)
		if err != nil {
			return nil, err
		}
		err = p.expect(TokenIdent)
		if err != nil {
			return nil, err
		}
		keyType := p.currentToken.Value
		err = p.advance()
		if err != nil {
			return nil, err
		}
		err = p.expectAndAdvance(TokenSymbol, constant.SymbolComma)
		if err != nil {
			return nil, err
		}
		err = p.expect(TokenIdent)
		if err != nil {
			return nil, err
		}
		valueType := p.currentToken.Value
		err = p.advance()
		if err != nil {
			return nil, err
		}
		field.MapInfo = &MapInfo{
			KeyType:   keyType,
			ValueType: valueType,
		}
		mapMessage := generateMapMessage(keyType, valueType)
		if mapMessage != nil {
			p.protoc.Messages = append(p.protoc.Messages, mapMessage)
		}
		err = p.expect(TokenSymbol, constant.SymbolGreaterThan)
		if err != nil {
			return nil, err
		}
		field.Repeated = true
		p.currentToken.Value = keyType + "_" + valueType + "_map_entry"
	}
	if p.currentToken.Value == constant.KeywordRepeated {
		field.Repeated = true
		if err := p.advance(); err != nil {
			return nil, err
		}
	}

	// 解析类型
	field.TypeName = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过类型名
		return nil, err
	}

	// 解析字段名
	if err := p.expect(TokenIdent); err != nil {
		return nil, err
	}
	field.Name = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过字段名
		return nil, err
	}

	// 解析标签号：'= 42'
	if err := p.expect(TokenSymbol, "="); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '='
		return nil, err
	}
	if err := p.expect(TokenNumber); err != nil {
		return nil, err
	}
	_, _ = fmt.Sscanf(p.currentToken.Value, "%d", &field.FieldNumber)
	if err := p.advance(); err != nil { // 跳过数字
		return nil, err
	}

	// 解析选项（如 '[Deprecated = true]'）
	if p.currentToken.Value == constant.SymbolLeftBracket {
		if err := p.parseFieldOptions(field); err != nil {
			return nil, err
		}
	}

	return field, p.expectAndAdvance(TokenSymbol, ";")
}

func generateMapMessage(keyType string, valueType string) *Message {
	return &Message{
		Name: keyType + "_" + valueType + "_map_entry",
		Fields: []*Field{
			{Name: "key", TypeName: keyType, FieldNumber: 1, WireType: str2WireType(keyType)},
			{Name: "value", TypeName: valueType, FieldNumber: 2, WireType: str2WireType(valueType)},
		},
	}
}

func (p *Parser) parseFieldOptions(field *Field) error {
	if err := p.advance(); err != nil { // 跳过 '['
		return err
	}

	for p.currentToken.Value != constant.SymbolRightBracket {
		// 解析选项名（如 "Deprecated"）
		optionName := p.currentToken.Value
		if err := p.advance(); err != nil {
			return err
		}

		// 解析 '='
		if err := p.expect(TokenSymbol, "="); err != nil {
			return err
		}
		if err := p.advance(); err != nil {
			return err
		}

		// 解析选项值（支持布尔值、字符串、数字）
		var optionValue interface{}
		switch p.currentToken.Type {
		case TokenIdent:
			switch p.currentToken.Value {
			case constant.DefaultTrue:
				optionValue = true
			case constant.DefaultFalse:
				optionValue = false
			default:
				optionValue = p.currentToken.Value // 自定义选项（如 "foo.bar"）
			}
		case TokenString:
			optionValue = p.currentToken.Value
		case TokenNumber:
			num, _ := strconv.Atoi(p.currentToken.Value)
			optionValue = num
		default:
			return fmt.Errorf("invalid option value: %v", p.currentToken.Value)
		}

		// 记录选项（此处简化为处理已知选项）
		switch optionName {
		case constant.OptionDeprecated:
			field.Options.Deprecated = optionValue.(bool)
		case constant.OptionPacked:
			field.Options.Packed = optionValue.(bool)
		}

		if err := p.advance(); err != nil {
			return err
		}

		// 跳过可能的逗号分隔符
		if p.currentToken.Value == constant.SymbolComma {
			if err := p.advance(); err != nil {
				return err
			}
		}
	}

	return p.advance() // 跳过 ']'
}

func (p *Parser) parseOneOf() (*OneOf, error) {
	oneof := &OneOf{}
	if err := p.advance(); err != nil { // 跳过 'oneof'
		return nil, err
	}
	if err := p.expect(TokenIdent); err != nil {
		return nil, err
	}
	oneof.Name = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过名称
		return nil, err
	}
	if err := p.expect(TokenSymbol, constant.SymbolLeftBrace); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '{'
		return nil, err
	}

	for p.currentToken.Value != constant.SymbolRightBrace {
		field, err := p.parseField()
		if err != nil {
			return nil, err
		}
		oneof.Fields = append(oneof.Fields, field)
	}

	if err := p.advance(); err != nil { // 跳过 '}'
		return nil, err
	}
	return oneof, nil
}

func (p *Parser) parseEnum() (*Enum, error) {
	enum := &Enum{}
	if err := p.advance(); err != nil { // 跳过 'enum'
		return nil, err
	}
	if err := p.expect(TokenIdent); err != nil {
		return nil, err
	}
	enum.Name = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过枚举名
		return nil, err
	}
	if err := p.expect(TokenSymbol, constant.SymbolLeftBrace); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '{'
		return nil, err
	}

	for p.currentToken.Value != constant.SymbolRightBrace {
		if p.currentToken.Value == constant.SymbolSemicolon {
			if err := p.advance(); err != nil {
				return nil, err
			}
			continue
		}
		value := &EnumValue{}
		value.Name = p.currentToken.Value
		if err := p.advance(); err != nil { // 跳过枚举项名
			return nil, err
		}
		if err := p.expect(TokenSymbol, constant.SymbolEqual); err != nil {
			return nil, err
		}
		if err := p.advance(); err != nil { // 跳过 '='
			return nil, err
		}
		if err := p.expect(TokenNumber); err != nil {
			return nil, err
		}
		_, _ = fmt.Sscanf(p.currentToken.Value, "%d", &value.Value)
		enum.Values = append(enum.Values, value)
		if err := p.advance(); err != nil { // 跳过数字
			return nil, err
		}
		if err := p.expect(TokenSymbol, constant.SymbolSemicolon); err != nil {
			return nil, err
		}
		if err := p.advance(); err != nil { // 跳过 ';'
			return nil, err
		}
	}

	if err := p.advance(); err != nil { // 跳过 '}'
		return nil, err
	}
	return enum, nil
}

func (p *Parser) parseService() (*Service, error) {
	service := &Service{}
	if err := p.advance(); err != nil { // 跳过 'service'
		return nil, err
	}
	if err := p.expect(TokenIdent); err != nil {
		return nil, err
	}
	service.Name = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过服务名
		return nil, err
	}
	if err := p.expect(TokenSymbol, "{"); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '{'
		return nil, err
	}

	for p.currentToken.Value != "}" {
		method, err := p.parseMethod()
		if err != nil {
			return nil, err
		}
		service.Methods = append(service.Methods, method)
	}

	if err := p.advance(); err != nil { // 跳过 '}'
		return nil, err
	}
	return service, nil
}

func (p *Parser) parseMethod() (*Method, error) {
	method := &Method{}
	if err := p.advance(); err != nil { // 跳过 'rpc'
		return nil, err
	}
	if err := p.expect(TokenIdent); err != nil {
		return nil, err
	}
	method.Name = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过方法名
		return nil, err
	}
	if err := p.expect(TokenSymbol, "("); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '('
		return nil, err
	}

	// 处理客户端流（输入类型前的 "stream"）
	if p.currentToken.Value == constant.KeywordStream {
		method.ClientStreaming = true
		if err := p.advance(); err != nil { // 跳过 "stream"
			return nil, err
		}
	}

	// 解析输入类型（如 "RequestType"）
	method.InputType = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过输入类型
		return nil, err
	}
	if err := p.expect(TokenSymbol, ")"); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 ')'
		return nil, err
	}

	// 解析返回类型
	if err := p.expect(TokenIdent, constant.KeywordReturns); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 "returns"
		return nil, err
	}
	if err := p.expect(TokenSymbol, "("); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 '('
		return nil, err
	}

	// 处理服务端流（输出类型前的 "stream"）
	if p.currentToken.Value == constant.KeywordStream {
		method.ServerStreaming = true
		if err := p.advance(); err != nil { // 跳过 "stream"
			return nil, err
		}
	}

	method.OutputType = p.currentToken.Value
	if err := p.advance(); err != nil { // 跳过输出类型
		return nil, err
	}
	if err := p.expect(TokenSymbol, ")"); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil { // 跳过 ')'
		return nil, err
	}

	return method, p.expectAndAdvance(TokenSymbol, ";")
}
