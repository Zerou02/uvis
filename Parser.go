package main

type NodeType int

const (
	None NodeType = iota + 1
	Package
	Import
	ClassDecl
	Decl
	Assignment
	FunctionHeader
)

type Node struct {
	nodeType NodeType
	value    string
}

type Parser struct {
	tokens   *[]string
	nodes    []Node
	currIdx  int
	tokenLen int
}

func newParser() Parser {
	return Parser{}
}

func (this *Parser) parse(tokens *[]string) {
	this.tokens = tokens
	this.nodes = []Node{}
	this.currIdx = 0
	this.tokenLen = len(*tokens)

	for this.currIdx < this.tokenLen {
		var currToken = (*this.tokens)[this.currIdx]
		if currToken == "package" {
			this.parsePackage()
		} else if currToken == "import" {
			this.parseImport()
		} else if currToken == "class" {
			this.parseClassDecl()
		} else {

			this.currIdx++
		}
	}
}
func (this *Parser) consumeToken() {
	this.currIdx++
}

func (this *Parser) assertAndConsume(expected string) {
	if this.getCurrToken() != expected {
		for _, x := range this.nodes {
			println(x.value)
		}
		panic("Expected" + expected + ",but got " + this.getCurrToken())
	} else {
		this.currIdx++
	}
}
func (this *Parser) getCurrToken() string {
	return (*this.tokens)[this.currIdx]
}

func (this *Parser) parsePackage() {
	this.currIdx++
	var newNode = Node{
		nodeType: Package,
		value:    this.getCurrToken(),
	}
	this.currIdx++
	this.assertAndConsume(";")
	this.nodes = append(this.nodes, newNode)
}

func (this *Parser) parseImport() {
	this.currIdx++
	var newNode = Node{
		nodeType: Import,
		value:    this.getCurrToken(),
	}
	this.currIdx++
	this.assertAndConsume(";")
	this.nodes = append(this.nodes, newNode)
}

func (this *Parser) isVisModifier(token string) bool {
	return token == "public" || token == "private" || token == "protected"
}

func (this *Parser) isStaticModifier(token string) bool {
	return token == "static"
}

func (this *Parser) parseClassDecl() {
	this.currIdx -= 2
	var nodeVal = ""
	if this.isVisModifier(this.getCurrToken()) {
		nodeVal += this.getCurrToken() + " "
	}
	this.currIdx++
	if this.isVisModifier(this.getCurrToken()) || this.isStaticModifier(this.getCurrToken()) {
		nodeVal += this.getCurrToken() + " "
	}
	this.currIdx++
	nodeVal += this.getCurrToken() + " "
	this.currIdx++
	for _, x := range this.getTokensUntil("{") {
		nodeVal += x + " "
	}
	this.nodes = append(this.nodes, Node{
		nodeType: ClassDecl,
		value:    nodeVal,
	})

	//jetzt entweder assignment,declaration oder function
	for this.currIdx < this.tokenLen {
		var next = this.getTokensUntilEither([]string{"=", "(", ";"})
		var funcStr = "("
		var assgnmStr = "="
		if this.contains(&next, &funcStr) {
			this.parseFunc(&next)
		} else if this.contains(&next, &assgnmStr) {
			this.parseAssignment(&next)
		} else {
			this.parseDecl()
		}
	}
}

func (this *Parser) parseDecl() {
	//Lesekopf auf Zeichen hinter schließendem Semikolon
	this.currIdx -= 5
	var newTokens = ""
	if this.isVisModifier(this.getCurrToken()) {
		newTokens += this.getCurrToken() + " "
	}
	this.currIdx++
	if this.isStaticModifier(this.getCurrToken()) || this.isVisModifier(this.getCurrToken()) {
		newTokens += this.getCurrToken() + " "
	}
	this.currIdx++
	newTokens += this.getCurrToken() + " "
	this.currIdx++
	newTokens += this.getCurrToken() + " "
	this.currIdx++

	this.nodes = append(this.nodes, Node{
		nodeType: Decl,
		value:    newTokens,
	})
	this.currIdx++
}

func (this *Parser) isUppercase(str string) bool {
	var first = str[0]
	return first >= 'A' && first <= 'Z'
}
func (this *Parser) parseFunc(tokensTillNow *[]string) {
	this.currIdx--
	var newVal = ""
	var isConstructor = this.isUppercase(this.getCurrToken())
	if isConstructor {
		this.currIdx--
	} else {
		this.currIdx -= 2
	}
	if this.isVisModifier(this.getCurrToken()) {
		newVal += this.getCurrToken() + " "
	}
	this.currIdx++
	newVal += this.getCurrToken() + " "
	this.currIdx++
	newVal += this.getCurrToken() + " "
	this.assertAndConsume("(")
	for this.getCurrToken() != ")" {
		newVal += this.getCurrToken() + " "
		this.currIdx++
	}
	newVal += this.getCurrToken()
	this.currIdx++
	this.nodes = append(this.nodes, Node{
		nodeType: FunctionHeader,
		value:    newVal,
	})
	this.discardBody()
}

func (this *Parser) discardBody() {
	this.assertAndConsume("{")
	var closedParens = 0
	var openParens = 1
	for closedParens != openParens {
		var c = this.getCurrToken()
		if c == "{" {
			openParens++
		} else if c == "}" {
			closedParens++
		}
		this.currIdx++
	}
	this.currIdx++
}

func (this *Parser) parseAssignment(tokensTillNow *[]string) {
	this.parseDecl()
	this.getTokensUntil(";")
}

func (this *Parser) getTokensUntil(token string) []string {
	return this.getTokensUntilEither([]string{token})
}

func (this *Parser) contains(arr *[]string, x *string) bool {
	for _, y := range *arr {
		if *x == y {
			return true
		}
	}
	return false
}
func (this *Parser) backTrackUntilEither(tokens []string){
	for !this.contains()
}

// consumes searched for token too
func (this *Parser) getTokensUntilEither(tokens []string) []string {
	var retVal = []string{}
	var curr = this.getCurrToken()
	retVal = append(retVal, curr)
	for !this.contains(&tokens, &curr) && this.currIdx < this.tokenLen-1 {
		this.currIdx++
		curr = this.getCurrToken()
		retVal = append(retVal, curr)
	}
	this.currIdx++
	return retVal
}
