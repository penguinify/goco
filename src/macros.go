package src

import (
	"fmt"
	"os"
	"strconv"
	"unicode"

)

func GetMacroList(path string) []string {
	files, _ := os.ReadDir(path)
	var macros []string

	for _, file := range files {
		macros = append(macros, file.Name())
	}

	return macros
}


type TokenType int

const (
    TokenFunction TokenType = iota
    TokenKeyword
    TokenString
    TokenNumber
)

var Keywords = [...]string{"loop", "end", "forever"}

type Token struct {
    Type TokenType
    Value string
}

type Parser struct {
    Tokens []Token
    pos int
}

type Lexer struct {
    Source string
    pos int
}

type ASTNode struct {
    Type TokenType
    Value string
    Children []*ASTNode
}

func NewParser(input string) *Parser {
    lexer := NewLexer(input)
    tokens := lexer.Tokenize()
    return &Parser{Tokens: tokens}
}


func NewLexer(source string) *Lexer {
    return &Lexer{Source: source}
}

func (l *Lexer) Tokenize() []Token {
    var tokens []Token

    for l.pos < len(l.Source) {

        char := l.Source[l.pos]

        switch {
        case unicode.IsSpace(rune(char)):
            l.pos++
            continue
        case char == '"':
            tokens = append(tokens, Token{Type: TokenString, Value: l.String()})
        case unicode.IsDigit(rune(char)):
            tokens = append(tokens, Token{Type: TokenNumber, Value: strconv.Itoa(l.Number())})
        default:
            value := l.readWord()

            if l.IsKeyword(value) {
                tokens = append(tokens, Token{Type: TokenKeyword, Value: value})
            } else {
                tokens = append(tokens, Token{Type: TokenFunction, Value: value})
            }

        }
        l.pos++
    }


    return tokens
}


func (l *Lexer) IsKeyword(keyword string) bool {
    for _, k := range Keywords {
        if k == keyword {
            return true
        }
    }

    return false
}

func (l *Lexer) readWord() string {
    var keyword string

    for l.pos < len(l.Source) {
        if unicode.IsSpace(rune(l.Source[l.pos])) {
            break
        }

        keyword += string(l.Source[l.pos])
        l.pos++
    }

    return keyword
}

func (l *Lexer) String() string {
    var str string

    l.pos++

    for l.pos < len(l.Source) {
        if l.Source[l.pos] == '"' {
            break
        }

        str += string(l.Source[l.pos])
        l.pos++
    }

    return str
}

func (l *Lexer) Number() int {
    var num string

    for l.pos < len(l.Source) {
        if unicode.IsSpace(rune(l.Source[l.pos])) {
            break
        }

        num += string(l.Source[l.pos])
        l.pos++
    }

    n, _ := strconv.Atoi(num)
    return n
}

func (p *Parser) Parse() *ASTNode {
    parent := &ASTNode{Type: 0, Value: "root", Children: []*ASTNode{}}

    // Example Syntax
    /*
    0 is Function
    1 is Keyword
    2 is String
    3 is Number
mouseset 10 10
click "left"
forever
loop 3
type "hello" 0.5
end
keypress "a"
keyrelease "a"
    */
    for p.pos < len(p.Tokens) {
        token := p.Tokens[p.pos]
        
        switch token.Type {
        case TokenFunction:
            parent.Children = append(parent.Children, &ASTNode{Type: TokenFunction, Value: token.Value})
            p.pos++
        case TokenKeyword:
            keywordNode := &ASTNode{Type: TokenKeyword, Value: token.Value, Children: []*ASTNode{}}
            parent.Children = append(parent.Children, keywordNode)
            p.pos++



            switch token.Value {
            case "loop":
                
                keywordNode.Children = append(parent.Children, &ASTNode{Type: TokenNumber, Value: p.Tokens[p.pos].Value})
                p.pos++

                keywordNode.Children = p.Parse().Children
                
                if p.pos >= len(p.Tokens) {
                    break
                }

            case "forever":
                // no end keyword will stop this, this goes to the end of the file
                for p.pos < len(p.Tokens) {
                    keywordNode.Children = p.Parse().Children
                }

            case "end":
                return parent
            }


        case TokenNumber:
            parent.Children = append(parent.Children, &ASTNode{Type: TokenNumber, Value: token.Value})
            p.pos++
        case TokenString:
            parent.Children = append(parent.Children, &ASTNode{Type: TokenString, Value: token.Value})
            p.pos++
        }
    }

    return parent
}

