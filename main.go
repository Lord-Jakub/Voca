package main

//imports
import (
	"Voca/lib"
	"Voca/num"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// Create interpreter structure.
type Interpret struct {
	tokens   []Token
	KeyWords []string
}

// Create code structure. There are stored vars
type code struct {
	vars  map[string]string
	ivars map[string]string
}
type TokenType int

const (
	Invalid     = iota
	OpenParen   //1
	CloseParen  //2
	OpenBrace   //3
	CloseBrace  //4
	Plus        //5
	Minus       //6
	Multiply    //7
	Divide      //8
	Backslash   //9
	NewLine     //10
	SingleQuote //11
	DoubleQuote //12
	Equal       //13
	Not         //14
	Comma       //15
	LessThan    //16
	MoreThan    //17
	Int         //18
	String      //19
	Text        //20
	Keyword     //21
	Whitespace  //22
	DoubleEqual //23
	NotEqual    //24
)

// Token contains meta about each token
type Token struct {
	Type  TokenType
	Value any
	Line  int
}

var symbolMap = map[byte]TokenType{
	'(':  OpenParen,
	')':  CloseParen,
	'{':  OpenBrace,
	'}':  CloseBrace,
	'+':  Plus,
	'-':  Minus,
	'*':  Multiply,
	'/':  Divide,
	'\\': Backslash,
	'=':  Equal,
	'!':  Not,
	',':  Comma,
	'<':  LessThan,
	'>':  MoreThan,
}

// laxer function
func (i *Interpret) lexer(input string) {
	//create map of tokens
	i.tokens = []Token{}

	pos := 0
	tokpos := 0
	lines := 1
	for pos < len(input) {

		c := input[pos]
		switch {
		case unicode.IsDigit(rune(c)):
			//Numbers
			var num string
			for pos < len(input) && (unicode.IsDigit(rune(input[pos])) || string(input[pos]) == ".") {
				num += string(input[pos])
				pos++
			}
			i.tokens = append(i.tokens, Token{
				Type:  Int,
				Value: num,
				Line:  lines,
			})

			pos--
			tokpos++
		case unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) && !(string(c) == " ") && !(string(c) == "\"") && !(string(c) == "'"):
			//Strings
			var s string
			for pos < len(input) && (unicode.IsLetter(rune(input[pos])) || unicode.IsDigit(rune(input[pos])) || (string(input[pos]) == ".")) && !(string(input[pos]) == " ") && !(string(input[pos]) == "\"") && !(string(input[pos]) == "'") {

				s += string(input[pos])
				pos++
			}
			if lib.Contains(s, i.KeyWords) {
				i.tokens = append(i.tokens, Token{
					Type:  Keyword,
					Value: s,
					Line:  lines,
				})

				tokpos++
			} else {
				i.tokens = append(i.tokens, Token{
					Type:  String,
					Value: s,
					Line:  lines,
				})

				tokpos++
			}
			pos--

		case string(c) == "/" && string(input[pos+1]) == "/":
			//Comments
			for pos < len(input) && string(input[pos]) != "\n" {
				pos++
			}
			lines++
		case string(c) == "=" && string(input[pos+1]) == "=":
			pos++
			i.tokens = append(i.tokens, Token{
				Type: DoubleEqual,
				Line: lines,
			})

		case string(c) == "!" && string(input[pos+1]) == "=":
			pos++
			i.tokens = append(i.tokens, Token{
				Type: NotEqual,
				Line: lines,
			})

		case string(c) == "\"":
			//Strings
			var s string
			pos++
			for pos < len(input) && string(input[pos]) != "\"" {
				s += string(input[pos])
				pos++
			}

			i.tokens = append(i.tokens, Token{
				Type:  Text,
				Value: s,
				Line:  lines,
			})
		case string(c) == "'":
			//Strings
			var s string
			pos++
			for pos < len(input) && string(input[pos]) != "'" {
				s += string(input[pos])
				pos++
			}

			i.tokens = append(i.tokens, Token{
				Type:  Text,
				Value: s,
				Line:  lines,
			})
		case string(c) == "\n":
			//New line
			i.tokens = append(i.tokens, Token{
				Type: NewLine,
				Line: lines,
			})
			lines++
		default:
			if token, ok := symbolMap[c]; ok {
				i.tokens = append(i.tokens, Token{Type: token, Line: lines, Value: string(c)})
				tokpos++
			}
		}
		pos++
	}

}

// GetCode function to extract code between curly braces
func GetCode(tokens []Token, i int) ([]Token, int) {
	code := []Token{}
	for tokens[i].Type != OpenBrace {
		i++
	}
	if tokens[i].Type == OpenBrace {
		n := 1
		i++
		x := 0
		for n != 0 {
			if tokens[i].Type == OpenBrace {
				n++
			}
			if tokens[i].Type == CloseBrace {
				n--
			}
			if n == 0 {
				break
			}
			code = append(code, tokens[i])
			x++
			i++
		}
	}
	code = append(code, Token{Type: NewLine})
	return code, i
}

// get name of func or var
func getname(tokens []Token, i int) string {
	for tokens[i].Type != String {
		i++
	}
	return tokens[i].Value.(string)
}

// get value of some tokens
func getvalue(tokens []Token, i int, vars map[string]string, fun map[string][]Token) (string, int) {
	for (tokens[i].Type != NewLine) && (tokens[i].Type != Comma) && (tokens[i].Type != CloseParen) {
		num1 := i
		for tokens[num1].Type != NewLine && (tokens[i].Type != Comma) && (tokens[i].Type != CloseParen) {
			if tokens[i].Type == String {
				token := tokens[i]

				if value, exists := vars[token.Value.(string)]; exists {
					if floatValue, isFloat := strconv.ParseFloat(value, 64); isFloat == nil {
						tokens[i] = Token{Type: Int, Value: strconv.FormatFloat(floatValue, 'f', 6, 64)}
					}
				} else if token.Value == "true" || token.Value == "false" {
					return token.Value.(string), i

				} else if _, exists := fun[token.Value.(string)]; exists {
					// Get the function arguments and prepare for function execution
					i3 := i
					for tokens[i3].Type != OpenParen {
						i3++
					}
					i3++
					i4 := 0
					fargs := make(map[int]string)
					for tokens[i3].Type != CloseParen && tokens[i3].Type != NewLine {
						if tokens[i3].Type != Comma {
							//i4b := 0
							fargs[i4], _ = getvalue(tokens, i3, vars, fun)
							i4++ //= i4b
						}

						i3++
					}
					// Prepare for function execution
					fnum := 0
					fvars1 := make(map[string]string)
					funcp := make(map[string][]Token)
					for key, value := range fun {
						funcp[key] = value
					}
					funCopy := make(map[string][]Token)
					for k, v := range fun {
						funCopy[k] = []Token{}
						for _, v2 := range v {
							funCopy[k] = append(funCopy[k], v2)
						}
					}
					if funcp[tokens[i].Value.(string)][fnum].Type == String {
						for funcp[tokens[i].Value.(string)][fnum].Type == String && strings.HasPrefix(funcp[tokens[i].Value.(string)][fnum].Value.(string), "VAR:") {
							// Assign values to function parameters
							fvars1[strings.TrimPrefix(funcp[tokens[i].Value.(string)][fnum].Value.(string), "VAR:")] = fargs[fnum]
							funcp[tokens[i].Value.(string)][fnum].Type = Whitespace
							fnum++
						}
					}

					// Create a new code instance for function execution
					c2 := code{
						vars: fvars1,
					}
					fun = funCopy
					// Execute the function code

					val := c2.Code(funcp[tokens[i].Value.(string)], fun)

					if floatValue, isFloat := strconv.ParseFloat(val, 64); isFloat == nil {
						tokens[i] = Token{Type: Int, Value: strconv.FormatFloat(floatValue, 'f', 6, 64)}
					} else {
						tokens[i] = Token{Type: Text, Value: val}
					}
					int1 := i + 1
					for int1 < i3 {
						if tokens[int1].Type != Whitespace {
							tokens[int1].Type = Whitespace
						}
						int1++
					}
					i--
				}
			}
			num1++
		}
		if tokens[i].Type == String {
			if _, exist := vars[tokens[i].Value.(string)]; exist {
				if _, err := strconv.Atoi(vars[tokens[i].Value.(string)]); err == nil {
					tokens[i] = Token{Type: Int, Value: vars[tokens[i].Value.(string)]}
				}

			}
		}
		if tokens[i].Type == Text {
			s := tokens[i].Value.(string)
			joinstr := false
			i2 := i
			for tokens[i2].Type != NewLine {
				if tokens[i2].Type == Plus {
					joinstr = true
				}
				i2++
			}
			i++
			if joinstr {
				for tokens[i].Type != NewLine {
					if tokens[i].Type == Text {
						s += tokens[i].Value.(string)
					} else if tokens[i].Type == Int {
						s += tokens[i].Value.(string)
					} else if tokens[i].Type == String {
						// If the token is a variable, replace it with its value

						if _, exists := vars[tokens[i].Value.(string)]; exists {
							s += vars[tokens[i].Value.(string)]
						} else if tokens[i].Value.(string) == "Random" {
							min := 0
							max := 0
							for tokens[i].Type != OpenParen {
								i++
							}
							i++
							for tokens[i].Type != Comma {
								if tokens[i].Type == Int {

									//convert to int
									min, _ = strconv.Atoi(tokens[i].Value.(string))
								}
								i++
							}
							for tokens[i].Type != CloseParen {
								if tokens[i].Type == Int {

									//convert to int
									max, _ = strconv.Atoi(tokens[i].Value.(string))
								}
								i++
							}

							randomNumber := rand.Intn(max+1-min) + min

							s += strconv.Itoa(randomNumber)
						} else if tokens[i].Value.(string) == "Read" {
							i++
							i++
							s += lib.Read()
						}
					}
					i++
				}
			}
			// Extract the text content and return

			return s, i
		} else if tokens[i].Type == Int {
			// If the token is an integer
			numb := make(map[int]string)
			numbs := 0
			for tokens[i].Type != NewLine && tokens[i].Type != Comma && tokens[i].Type != CloseParen {
				// Check if the token is a variable and replace it with its value
				if tokens[i].Type == String {

					if _, exists := vars[tokens[i].Value.(string)]; exists {
						tokens[i] = Token{Type: Int, Value: vars[tokens[i].Value.(string)]}
					}
				}

				// Collect numeric tokens and operators
				if tokens[i].Type == Int || tokens[i].Type == Plus || tokens[i].Type == Minus || tokens[i].Type == Multiply || tokens[i].Type == Divide {
					if tokens[i].Type == Int {
						numb[numbs] = tokens[i].Value.(string)
					} else {
						if tokens[i].Type == Plus {
							numb[numbs] = "+"
						} else if tokens[i].Type == Minus {
							numb[numbs] = "-"
						} else if tokens[i].Type == Multiply {
							numb[numbs] = "*"
						} else if tokens[i].Type == Divide {
							numb[numbs] = "/"
						}
					}
					numbs++
				}

				i++
			}
			// Evaluate the numeric expression and return the result as a string
			res, _ := num.Evaluate(numb)
			return strconv.FormatFloat(res, 'f', 6, 64), i
		} else if tokens[i].Type == String {
			// If the token is a variable, replace it with its value

			if _, exists := vars[tokens[i].Value.(string)]; exists {
				return vars[tokens[i].Value.(string)], i
			} else if tokens[i].Value.(string) == "Random" {
				min := 0
				max := 0
				for tokens[i].Type != OpenParen {
					i++
				}
				i++
				for tokens[i].Type != Comma {
					if tokens[i].Type == Int {

						//convert to int
						min, _ = strconv.Atoi(tokens[i].Value.(string))
					}
					i++
				}
				for tokens[i].Type != CloseParen {
					if tokens[i].Type == Int {

						//convert to int
						max, _ = strconv.Atoi(tokens[i].Value.(string))
					}
					i++
				}

				randomNumber := rand.Intn(max+1-min) + min

				return strconv.Itoa(randomNumber), i
			} else if tokens[i].Value.(string) == "Read" {
				i++
				i++
				return lib.Read(), i
			}
		} else {
		}

		i++
	}
	if tokens[i].Type != NewLine && tokens[i-1].Type != Whitespace {
		lib.Print("Unexpected token: \"" + tokens[i].Value.(string) + "\" on line " + strconv.Itoa(tokens[i].Line))
	}
	return "", i
}

func Contains(s TokenType, array []TokenType) bool {
	for _, value := range array {
		if s == value {
			return true
		}
	}
	return false
}

// getbool function to evaluate boolean expressions
func getbool(tokens []Token, i int, vars map[string]string, fun map[string][]Token) bool {
	// List of operators
	ops := []TokenType{DoubleEqual, LessThan, MoreThan, Not}

	// Create a map for the first set of tokens
	toks1 := []Token{}
	i2 := 0

	// Collect tokens until an operator is found
	for !Contains(tokens[i].Type, ops) {
		toks1 = append(toks1, tokens[i])
		i++
		i2++
	}

	op := tokens[i].Type
	/*if op == DoubleEqual || op == NotEqual {
		i++
	}*/
	i++

	// Create a map for the second set of tokens
	toks2 := []Token{}
	i2 = 0

	// Collect tokens until a newline or opening brace is encountered
	for (tokens[i].Type != NewLine) && (tokens[i].Type != OpenBrace) {
		toks2 = append(toks2, tokens[i])
		i++
		i2++
	}

	// Add "NEWLINE" tokens to both sets
	toks1 = append(toks1, Token{Type: NewLine})
	toks2 = append(toks2, Token{Type: NewLine})

	// Get values from the token sets
	val1, _ := getvalue(toks1, 0, vars, fun)
	val2, _ := getvalue(toks2, 0, vars, fun)
	val1i, _ := strconv.Atoi(val1)
	val2i, _ := strconv.Atoi(val2)

	// Compare values based on the operator and return the result
	if (op == DoubleEqual) && (val1 == val2) {
		return true
	} else if (op == NotEqual) && (val1 != val2) {
		return true
	} else if (op == LessThan) && (val1i < val2i) {
		return true
	} else if (op == MoreThan) && (val1i > val2i) {
		return true
	}
	return false
}

// interpret function to process and execute the interpreted code
func interpret(tokens []Token) {

	// Initialize a map to store functions
	fun := make(map[string][]Token)

	// Initialize a map for instance variables
	ivars := make(map[string]string)

	// Initialize the index variable
	i := 0

	// Loop through the tokens
	for i < len(tokens) {
		// Initialize a map for the current function's code
		funcode := []Token{}

		// Check if the current token indicates the start of a function
		if tokens[i].Type == Keyword && tokens[i].Value == "func" {
			// Initialize variables for function code extraction
			i2 := 0
			i3 := i

			// Go to an opening brace
			for tokens[i3].Type != OpenParen {
				i3++
			}
			i3++
			i4 := 0
			for tokens[i3].Type != CloseParen {
				// This process func parameters
				if tokens[i3].Type == String {
					funcode[i4] = Token{Type: String, Value: "VAR:" + tokens[i3].Value.(string)}
					i4++
				}

				i3++
			}

			// Extract the function code using the GetCode function
			funcode2, i22 := GetCode(tokens, i)

			n := 0
			// Add the function code to the map
			for n < len(funcode2) {
				funcode = append(funcode, funcode2[n])
				n++
			}

			i2 += i22
			fname := getname(tokens, i)
			i = i2
			fun[fname] = funcode
		} else if tokens[i].Type == Keyword && tokens[i].Value == "import" {
			// If the token indicates an import statement
			// Create a new Interpret instance for the imported file
			in := Interpret{
				KeyWords: []string{"print", "if", "var", "func", "while", "import", "return"},
			}

			// Move to the next token until a file path is found
			for tokens[i].Type != Text {
				i++
			}
			filename := tokens[i].Value.(string)
			// Get the file path and read the content of the file
			executablePath, err := os.Executable()
			if err != nil {
				lib.Print("Nelze získat cestu k spustitelnému souboru:" + err.Error())
				return
			}
			file_path := filepath.Dir(executablePath) + "/" + "libs/" + tokens[i].Value.(string) + ".v"
			data, _ := os.ReadFile(file_path)

			// Replace line endings and tokenize the content
			input := string(data)
			input = strings.Replace(input, "\r\n", "\n", -1)
			in.lexer(input)

			// Loop through the tokens of the imported file
			n := 0
			for n < len(in.tokens) {
				// Check if the token indicates the start of a function in the imported file
				if in.tokens[n].Type == Keyword && in.tokens[n].Value == "func" {
					// Initialize variables for function code extraction
					i2 := 0
					i3 := n

					for in.tokens[i3].Type != OpenParen {
						i3++
					}
					i3++
					i4 := 0
					for in.tokens[i3].Type != CloseParen {

						if in.tokens[i3].Type == String {
							funcode = append(funcode, Token{Type: String, Value: "VAR:" + in.tokens[i3].Value.(string)})
							i4++
						}

						i3++
					}

					// Extract the function code using the GetCode function
					funcode2, i22 := GetCode(in.tokens, n)
					n2 := 0
					// Add the function code to the map
					for n2 < len(funcode2) {
						funcode = append(funcode, funcode2[n2])
						n2++
					}

					// Update indices and store the function code with its name
					i2 += i22
					fname := filename + "." + getname(in.tokens, i)
					n = i2
					fun[fname] = funcode
				}

				// Move to the next token in the imported file
				n++
			}
		}

		// Move to the next token in the original file
		i++
		// Check if i exceeds the length of tokens
		if i >= len(tokens) {
			break
		}
	}

	// Create a code instance with variable and instance variable maps
	c := code{
		vars:  make(map[string]string),
		ivars: ivars,
	}

	// Execute the code for the main function
	c.Code(fun["main"], fun)
}

// Code executes the code represented by the given tokens and function map.
// It loops through the tokens to find and store functions, then loops through the tokens again to execute the code.
// It handles keywords such as "print", "var", "if", and "while", as well as function calls.
// The function map stores the code for each function, indexed by function name.
// The vars map stores the values of variables, indexed by variable name.
func (c *code) Code(tokens []Token, fun map[string][]Token) string {
	// Initialize the index variable
	i := 0
	// Initialize a map for function code
	funcode := []Token{}

	// Loop through the tokens to find and store functions
	for i < len(tokens) {
		// Check if the current token indicates the start of a function
		if tokens[i].Type == Keyword && tokens[i].Value == "func" {
			// Initialize variables for function code extraction
			i2 := 0
			i3 := i

			for tokens[i3].Type != OpenParen {
				i3++
			}
			i3++
			i4 := 0
			for tokens[i3].Type != CloseParen {

				if tokens[i3].Type == String {
					funcode[i4].Value = "VAR:" + tokens[i3].Value.(string)
					i4++
				}

				i3++
			}

			// Extract the function code using the GetCode function
			funcode2, i22 := GetCode(tokens, i)
			n := 0
			// Add the function code to the map
			for n < len(funcode2) {
				funcode[len(funcode)] = funcode2[n]
				n++
			}

			// Update indices and store the function code with its name
			i2 += i22
			fname := getname(tokens, i)
			i = i2
			fun[fname] = funcode
		}

		i++
	}

	// Reset index variable
	i = 0

	// Loop through the tokens to execute the code
	for i < len(tokens) {
		switch {
		case tokens[i].Type == Whitespace:
			// Ignore whitespace tokens

		case tokens[i].Type == NewLine:
			//Increase lines counter

		case tokens[i].Type == Keyword:
			// Check for keyword tokens
			switch {
			case tokens[i].Value == "print":
				// If the keyword is "print," get the value and print it
				val, vl := getvalue(tokens, i, c.vars, fun)
				lib.Print(val)
				i = vl
			case tokens[i].Value == "var":
				// If the keyword is "var," get the variable name and value, then store it
				fname := getname(tokens, i)
				for tokens[i].Type != Equal {
					i++
				}
				i++
				val, vl := getvalue(tokens, i, c.vars, fun)
				c.vars[fname] = val
				i = vl
			case tokens[i].Value == "if":
				// If the keyword is "if," get the code block and execute it if the condition is true
				i++
				toks, ifl := GetCode(tokens, i)

				if getbool(tokens, i, c.vars, fun) {
					c.Code(toks, fun)
				}
				i = ifl
			case tokens[i].Value == "while":
				// If the keyword is "while," get the code block and execute it while the condition is true
				i++
				toks, ifl := GetCode(tokens, i)

				for getbool(tokens, i, c.vars, fun) {
					toks, ifl = GetCode(tokens, i)

					c.Code(toks, fun)
				}
				i = ifl
			case tokens[i].Value == "return":
				// If the keyword is "return," get the value and return it
				i++
				val, _ := getvalue(tokens, i, c.vars, fun)
				return val

			}
		case tokens[i].Type == String:
			// If the token is a string, check if it corresponds to a function

			if _, exists := fun[tokens[i].Value.(string)]; exists {
				// Get the function arguments and prepare for function execution
				i3 := i
				for tokens[i3].Type != OpenParen {
					i3++
				}
				i3++
				i4 := 0
				fargs := make(map[int]string)
				for tokens[i3].Type != CloseParen && tokens[i3].Type != NewLine {
					if tokens[i3].Type != Comma {
						//i4b := 0
						fargs[i4], _ = getvalue(tokens, i3, c.vars, fun)
						i4++ //= i4b
					}

					i3++
				}
				// Prepare for function execution
				fnum := 0
				fvars1 := make(map[string]string)
				funcp := make(map[string][]Token)
				for key, value := range fun {
					funcp[key] = value
				}
				funCopy := make(map[string][]Token)
				for k, v := range fun {
					funCopy[k] = []Token{}
					for k2, v2 := range v {
						funCopy[k][k2] = v2
					}
				}
				for strings.HasPrefix(funcp[tokens[i].Value.(string)][fnum].Value.(string), "VAR:") {
					// Assign values to function parameters
					fvars1[strings.TrimPrefix(funcp[tokens[i].Value.(string)][fnum].Value.(string), "VAR:")] = fargs[fnum]
					funcp[tokens[i].Value.(string)][fnum].Type = Whitespace
					fnum++
				}

				// Create a new code instance for function execution
				c2 := code{
					vars: fvars1,
				}
				fun = funCopy
				// Execute the function code
				c2.Code(funcp[tokens[i].Value.(string)], fun)
				i = i3
			} else if _, exists := c.vars[tokens[i].Value.(string)]; exists {
				// If the token is a variable, replace it with its value
				fname := getname(tokens, i)
				for tokens[i].Type != Equal && i < len(tokens) {
					i++
				}
				i++
				val, vl := getvalue(tokens, i, c.vars, fun)
				c.vars[fname] = val
				i = vl
			} else {
				lib.Print("Unexpected token: \"" + tokens[i].Value.(string) + "\" on line " + strconv.Itoa(tokens[i].Line))
			}

		default:
			// Handle unknown keywords
			if tokens[i].Type == Whitespace {

			} else if tokens[i].Type == NewLine {
			} else {
				lib.Print("Unexpected token: \"" + tokens[i].Value.(string) + "\" on line " + strconv.Itoa(tokens[i].Line))
			}
		}
		// Move to the next token
		i++
	}
	return ""
}

func main() {

	i := Interpret{
		tokens:   make([]Token, 0),
		KeyWords: []string{"print", "if", "var", "func", "while", "import", "return"},
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "get" {
			//get if args exist
			if len(os.Args) > 2 {
				//get file path
				executablePath, _ := os.Executable()
				file_name, _ := lib.ExtractFileName(os.Args[2])
				if _, err := os.Stat(filepath.Dir(executablePath) + "/" + "libs"); os.IsNotExist(err) {
					os.Mkdir(filepath.Dir(executablePath)+"/"+"libs", 0755)
				}
				file_path := filepath.Dir(executablePath) + "/" + "libs/" + file_name
				//get urlr
				url := os.Args[2]
				//download file
				lib.Download(file_path, url)
			} else {
				lib.Print("Nesprávné použití příkazu get")
			}
		} else if os.Args[1] == "help" {
			lib.Print("Voca - programming language")
			lib.Print("Usage: -voca - to run main.v file")
			lib.Print("       -voca [file] - to run [file].v file")
			lib.Print("       -voca get [url] - to download file from [url] and save it to libs folder")
			lib.Print("       -voca help - to show this help")

		} else if os.Args[1] == "run" {
			//get file path
			file_path := os.Args[2]
			//read file
			data, _ := os.ReadFile(file_path)

			input := string(data)
			input = strings.Replace(input, "\\r\\n", "\\n", -1)
			i.lexer(input)
			interpret(i.tokens)
		} else if os.Args[1] == "version" {
			lib.Print("Voca version 0.1.2")
		}

	} else {
		file_path := ""
		//run if args exist

		cur_dir, _ := os.Getwd()
		file_path = cur_dir + "/main.v"

		data, _ := os.ReadFile(file_path)

		input := string(data)
		input = strings.Replace(input, "\\r\\n", "\\n", -1)
		i.lexer(input)
		interpret(i.tokens)
	}
}
