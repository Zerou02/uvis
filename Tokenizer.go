package main

type Tokenizer struct {
	str     *string
	length  int
	currIdx int
	tokens  []string
}

func newTokenizer() Tokenizer {
	return Tokenizer{
		str:     nil,
		length:  0,
		currIdx: 0,
		tokens:  []string{},
	}
}

func (this *Tokenizer) isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '.'
}
func (this *Tokenizer) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (this *Tokenizer) isAlphaNumerical(c rune) bool {
	return this.isAlpha(c) || this.isDigit(c)
}
func (this *Tokenizer) tokenizeAlphaNumerical() {
	var c = (*this.str)[this.currIdx]
	var token = ""
	for this.isAlphaNumerical(rune(c)) {
		token += string(c)
		this.currIdx++
		c = (*this.str)[this.currIdx]
	}
	this.tokens = append(this.tokens, token)
}

func (this *Tokenizer) tokenize(str *string) {
	this.str = str
	this.tokens = []string{}
	this.length = len(*str)
	this.currIdx = 0
	for this.currIdx < this.length {

		var c = (*this.str)[this.currIdx]
		if c == ' ' || c == '\n' || c == '\r' {
			//do nothing
			this.currIdx++
		} else if c == '(' || c == ')' || c == '{' || c == '}' || c == '=' || c == '<' || c == '>' || c == '.' || c == ';' || c == '[' || c == ']' || c == ':' || c == '+' || c == '-' || c == '*' || c == '/' || c == '@' || c == '"' || c == '\'' || c == '_' || c == '&' || c == '|' || c == ',' {
			this.tokens = append(this.tokens, string(c))
			this.currIdx++
		} else if this.isAlphaNumerical(rune(c)) {
			this.tokenizeAlphaNumerical()
		} else {
			println(string(rune(c)))
			panic("unrecognized token")
		}
	}
}
