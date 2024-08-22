package main

import (
	"io/fs"
	"os"
	"strings"
)

func main() {
	var args = os.Args
	if len(args) < 2 {
		println("invalid usage")
		return
	}
	var dir, err = os.ReadDir(args[1])
	if err != nil {
		panic(err.Error())
	}
	var files = []string{}
	collectFilesRecursively(args[1], dir, &files)

	var tokenizer = newTokenizer()
	for i := 0; i < len(files); i++ {
		var f, err = os.ReadFile(files[i])
		if err != nil {
			println(err)
		}
		var str = string(f)
		tokenizer.tokenize(&str)
	}

	println(files[1])
	var f, _ = os.ReadFile(files[1])
	var str = string(f)
	tokenizer.tokenize(&str)
	var parser = newParser()
	parser.parse(&tokenizer.tokens)

	for _, x := range parser.nodes {
		println(x.value)
	}
	var interpreter = newInterpreter()
	interpreter.start()
	interpreter.interpret(&parser.nodes)
	interpreter.end()

}

func collectFilesRecursively(path string, dir []fs.DirEntry, acc *[]string) {
	for _, x := range dir {
		if x.IsDir() {
			var subDir, err = os.ReadDir(path + x.Name() + "\\")
			if err != nil {
				panic(err.Error())
			}
			collectFilesRecursively(path+x.Name()+"\\", subDir, acc)
		} else {
			var splitted = strings.Split(x.Name(), ".")
			if len(splitted) == 2 && splitted[1] == "java" {
				*acc = append(*acc, path+x.Name())
			}
		}
	}
}

type Interpreter struct {
	content strings.Builder
}

func newInterpreter() Interpreter {
	return Interpreter{content: strings.Builder{}}
}

func (this *Interpreter) start() {
	this.content.WriteString("@startuml\n")
	this.content.WriteString("skinparam classAttributeIconSize 0\n")
}

func (this *Interpreter) interpret(nodes *[]Node) {
	for _, x := range *nodes {
		if x.nodeType == Package {
			this.interpretPackage(x.value)
		} else if x.nodeType == ClassDecl {
			this.interpretClassDecl(x.value)
		} else if x.nodeType == FunctionHeader {
			this.interpretFunctionHeader(x.value)
		} else if x.nodeType == Decl {
			this.interpretDecl(x.value)
		}
	}
	this.content.WriteString("}\n}\n") //Klasse und Package, fucking plantuml braucht 2 Zeilen
}

func (this *Interpreter) interpretDecl(value string) {
	var splitted = strings.Split(value, " ")
	var swap = splitted[0]
	splitted[0] = splitted[1] + ": "
	splitted[1] = swap
	this.content.WriteString(strings.Join(splitted, "") + "\n")
}

func (this *Interpreter) interpretFunctionHeader(value string) {
	var splitted = strings.Split(value, "(")
	splitted[0] += "("
	var furtherSplitted = strings.Split(splitted[1], " ")
	var new = []string{}
	for _, x := range furtherSplitted {
		if x != " " && x != "," && x != "" {
			new = append(new, x)
		}
	}
	furtherSplitted = new
	for i := 0; i < len(furtherSplitted)-1; i += 2 {
		var swapped = furtherSplitted[i]
		furtherSplitted[i] = furtherSplitted[i+1] + ":"
		furtherSplitted[i+1] = swapped
		if i < len(furtherSplitted)-3 {
			furtherSplitted[i+1] += ","
		}
	}
	splitted[0] = strings.ReplaceAll(splitted[0], "public", "+")
	splitted[0] = strings.ReplaceAll(splitted[0], "private", "-")
	splitted[0] = strings.ReplaceAll(splitted[0], "static", "{static}")
	splitted[0] = strings.ReplaceAll(splitted[0], "abstract", "{abstract}")
	var parens = strings.Join(furtherSplitted, " ")
	parens = strings.ReplaceAll(parens, " )", ")")
	this.content.WriteString(splitted[0] + parens + "\n")
}

func (this *Interpreter) interpretClassDecl(val string) {
	var splitted = strings.Split(strings.Split(val, "class")[1], " ")
	for i := 0; i < len(splitted); i++ {
		if splitted[i] == "<" {
			splitted[i] = "< ? extends"
		}
	}
	this.content.WriteString("class" + strings.Join(splitted, " "))
	this.content.WriteString("\n")
}

func (this *Interpreter) interpretPackage(val string) {
	this.content.WriteString("package " + strings.TrimPrefix(val, "com.") + "{\n")
}

func (this *Interpreter) end() {
	this.content.WriteString("@enduml\n")
	var f, err = os.Create("diagram.txt")
	if err != nil {
		println(err)
	}
	f.WriteString(this.content.String())
	f.Close()
}

func (this *Interpreter) contains(arr *[]string, x *string) bool {
	for _, y := range *arr {
		if *x == y {
			return true
		}
	}
	return false
}
