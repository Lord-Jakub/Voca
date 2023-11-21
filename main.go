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
	tokens   []string
	KeyWords []string
}

// Create code structure. There are stored vars
type code struct {
	vars  map[string]string
	ivars map[string]string
}

// laxer function
func (i *Interpret) lexer(input string) {
	//create map of tokens
	i.tokens = []string{}

	pos := 0
	tokpos := 0

	for pos < len(input) {
		c := input[pos]
		if unicode.IsDigit(rune(c)) {
			//Numbers
			var num string
			for pos < len(input) && (unicode.IsDigit(rune(input[pos])) || string(input[pos]) == ".") {
				num += string(input[pos])
				pos++
			}
			i.tokens = append(i.tokens, "INT:"+num)

			pos--
			tokpos++
		} else if unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) && !(string(c) == " ") {
			//Strings
			var s string
			for pos < len(input) && (unicode.IsLetter(rune(input[pos])) || unicode.IsDigit(rune(input[pos])) || (string(input[pos]) == ".")) && !(string(input[pos]) == " ") {

				s += string(input[pos])
				pos++
			}
			if lib.Contains(s, i.KeyWords) {
				i.tokens = append(i.tokens, "KEYWORD:"+s)

				tokpos++
			} else {
				i.tokens = append(i.tokens, "STRING:"+s)

				tokpos++
			}
			pos--

		} else if string(c) == " " {
			//Whitespaces
			i.tokens = append(i.tokens, "WHITESPACE")

			tokpos++
		} else {
			//Symbols
			if string(c) == "(" {
				i.tokens = append(i.tokens, "OP_B")
				tokpos++
			} else if string(c) == ")" {
				i.tokens = append(i.tokens, "CL_B")
				tokpos++
			} else if string(c) == "{" {
				i.tokens = append(i.tokens, "OP_S_B")
				tokpos++
			} else if string(c) == "}" {
				i.tokens = append(i.tokens, "CL_S_B")
				tokpos++
			} else if string(c) == "+" {
				i.tokens = append(i.tokens, "PLUS")
				tokpos++
			} else if string(c) == "-" {
				i.tokens = append(i.tokens, "MINUS")
				tokpos++
			} else if string(c) == "*" {
				i.tokens = append(i.tokens, "MULT")
				tokpos++
			} else if string(c) == "/" && string(input[pos+1]) != "/" {
				i.tokens = append(i.tokens, "DIV")
				tokpos++
			} else if string(c) == "\\" {
				i.tokens = append(i.tokens, "BACKSLASH")
				tokpos++
			} else if string(c) == "\n" {
				i.tokens = append(i.tokens, "NEWLINE")
				tokpos++
			} else if string(c) == "\"" {
				pos++
				var s string
				for string(input[pos]) != "\"" {
					s += string(input[pos])
					pos++
				}
				i.tokens = append(i.tokens, "TEXT:"+s)
				tokpos++
			} else if string(c) == "'" {
				pos++
				var s string
				for string(input[pos]) != "'" {
					s += string(input[pos])
					pos++
				}
				i.tokens = append(i.tokens, "TEXT:"+s)
				tokpos++
			} else if string(c) == "=" {
				if string(input[pos+1]) != "=" {
					i.tokens = append(i.tokens, "EQ")
				} else {
					i.tokens = append(i.tokens, "DEQ")
				}

				tokpos++
			} else if string(c) == "<" {
				i.tokens = append(i.tokens, "LESS")
				tokpos++
			} else if string(c) == ">" {
				i.tokens = append(i.tokens, "MORE")
				tokpos++
			} else if string(c) == "!" && string(input[pos+1]) == "=" {
				i.tokens = append(i.tokens, "NOT")
				tokpos++
			} else if string(c) == "," {
				i.tokens = append(i.tokens, "COM")
				tokpos++
			} else if string(c) == "/" && string(input[pos+1]) == "/" {
				for string(input[pos]) != "\n" {
					pos++
				}
				i.tokens = append(i.tokens, "NEWLINE")
				tokpos++
			}

		}
		pos++
	}

}

// GetCode function to extract code between curly braces
func GetCode(tokens []string, i int) ([]string, int) {
	code := []string{}
	for tokens[i] != "OP_S_B" {
		i++
	}
	if tokens[i] == "OP_S_B" {
		n := 1
		i++
		x := 0
		for n != 0 {

			if tokens[i] == "OP_S_B" {
				n++
			}
			if tokens[i] == "CL_S_B" {
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
	return code, i
}

// get name of func or var
func getname(tokens []string, i int) string {
	for !strings.HasPrefix(tokens[i], "STRING:") {
		i++
	}
	tokens[i] = strings.TrimPrefix(tokens[i], "STRING:")
	return tokens[i]
}

// get value of some tokens
func getvalue(tokens []string, i int, vars map[string]string, fun map[string][]string) (string, int) {
	for (tokens[i] != "NEWLINE") && (tokens[i] != "COM") && (tokens[i] != "CL_B") {
		num1 := i
		for tokens[num1] != "NEWLINE" && (tokens[i] != "COM") && (tokens[i] != "CL_B") {
			if strings.HasPrefix(tokens[i], "STRING:") {
				token := tokens[i]
				token = strings.TrimPrefix(token, "STRING:")

				if value, exists := vars[token]; exists {
					if floatValue, isFloat := strconv.ParseFloat(value, 64); isFloat == nil {
						tokens[i] = "INT:" + strconv.FormatFloat(floatValue, 'f', 6, 64)
					}
				} else if _, exists := fun[token]; exists {
					// Get the function arguments and prepare for function execution
					i3 := i
					for tokens[i3] != "OP_B" {
						i3++
					}
					i3++
					i4 := 0
					fargs := make(map[int]string)
					for tokens[i3] != "CL_B" && tokens[i3] != "NEWLINE" {
						if tokens[i3] != "COM" {
							//i4b := 0
							fargs[i4], _ = getvalue(tokens, i3, vars, fun)
							i4++ //= i4b
						}

						i3++
					}
					// Prepare for function execution
					fnum := 0
					fvars1 := make(map[string]string)
					funcp := make(map[string][]string)
					for key, value := range fun {
						funcp[key] = value
					}
					funCopy := make(map[string][]string)
					for k, v := range fun {
						funCopy[k] = []string{}
						for _, v2 := range v {
							funCopy[k] = append(funCopy[k], v2)
						}
					}
					tokens[i] = strings.TrimPrefix(tokens[i], "STRING:")
					for strings.HasPrefix(funcp[tokens[i]][fnum], "VAR:") {
						// Assign values to function parameters
						fvars1[strings.TrimPrefix(funcp[tokens[i]][fnum], "VAR:")] = fargs[fnum]
						funcp[tokens[i]][fnum] = "WHITESPACE"
						fnum++
					}

					// Create a new code instance for function execution
					c2 := code{
						vars: fvars1,
					}
					fun = funCopy
					// Execute the function code

					val := c2.Code(funcp[tokens[i]], fun)

					if floatValue, isFloat := strconv.ParseFloat(val, 64); isFloat == nil {
						tokens[i] = "INT:" + strconv.FormatFloat(floatValue, 'f', 6, 64)
					} else {
						tokens[i] = "TEXT:" + val
					}
					int1 := i + 1
					for int1 < i3 {
						if tokens[int1] != "WHITESPACE" {
							tokens[int1] = "WHITESPACE"
						}
						int1++
					}
					i--
				}
			}
			num1++
		}
		if _, exist := vars[tokens[i]]; exist {
			tokens[i] = "INT:" + vars[tokens[i]]
		}
		// Check if the token is a text string
		if strings.HasPrefix(tokens[i], "TEXT:") {
			s := strings.TrimPrefix(tokens[i], "TEXT:")
			joinstr := false
			i2 := i
			for tokens[i2] != "NEWLINE" {
				if tokens[i2] == "PLUS" {
					joinstr = true
				}
				i2++
			}
			i++
			if joinstr {
				for tokens[i] != "NEWLINE" {
					if strings.HasPrefix(tokens[i], "TEXT:") {
						s += strings.TrimPrefix(tokens[i], "TEXT:")
					} else if strings.HasPrefix(tokens[i], "INT:") {
						s += strings.TrimPrefix(tokens[i], "INT:")
					} else if strings.HasPrefix(tokens[i], "STRING:") {
						// If the token is a variable, replace it with its value
						tokens[i] = strings.TrimPrefix(tokens[i], "STRING:")
						if _, exists := vars[tokens[i]]; exists {
							s += vars[tokens[i]]
						} else if tokens[i] == "Random" {
							min := 0
							max := 0
							for tokens[i] != "OP_B" {
								i++
							}
							i++
							for tokens[i] != "COM" {
								if strings.HasPrefix(tokens[i], "INT:") {
									tokens[i] = strings.TrimPrefix(tokens[i], "INT:")
									//convert to int
									min, _ = strconv.Atoi(tokens[i])
								}
								i++
							}
							for tokens[i] != "CL_B" {
								if strings.HasPrefix(tokens[i], "INT:") {
									tokens[i] = strings.TrimPrefix(tokens[i], "INT:")
									//convert to int
									max, _ = strconv.Atoi(tokens[i])
								}
								i++
							}

							randomNumber := rand.Intn(max+1-min) + min

							s += strconv.Itoa(randomNumber)
						} else if tokens[i] == "Read" {
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
		} else if strings.HasPrefix(tokens[i], "INT:") {
			// If the token is an integer
			numb := make(map[int]string)
			numbs := 0
			for tokens[i] != "NEWLINE" && (tokens[i] != "COM") && (tokens[i] != "CL_B") {
				// Check if the token is a variable and replace it with its value
				if strings.HasPrefix(tokens[i], "STRING:") {
					tokens[i] = strings.TrimPrefix(tokens[i], "STRING:")
					if _, exists := vars[tokens[i]]; exists {
						tokens[i] = "INT:" + vars[tokens[i]]
					}
				}

				// Collect numeric tokens and operators
				if (strings.HasPrefix(tokens[i], "INT:")) || tokens[i] == "PLUS" || tokens[i] == "MINUS" || tokens[i] == "MULT" || tokens[i] == "DIV" {
					numb[numbs] = tokens[i]
					numbs++
				}

				i++
			}
			// Evaluate the numeric expression and return the result as a string
			res, _ := num.Evaluate(numb)
			return strconv.FormatFloat(res, 'f', 6, 64), i
		} else if strings.HasPrefix(tokens[i], "STRING:") {
			// If the token is a variable, replace it with its value
			tokens[i] = strings.TrimPrefix(tokens[i], "STRING:")
			if _, exists := vars[tokens[i]]; exists {
				return vars[tokens[i]], i
			} else if tokens[i] == "Random" {
				min := 0
				max := 0
				for tokens[i] != "OP_B" {
					i++
				}
				i++
				for tokens[i] != "COM" {
					if strings.HasPrefix(tokens[i], "INT:") {
						tokens[i] = strings.TrimPrefix(tokens[i], "INT:")
						//convert to int
						min, _ = strconv.Atoi(tokens[i])
					}
					i++
				}
				for tokens[i] != "CL_B" {
					if strings.HasPrefix(tokens[i], "INT:") {
						tokens[i] = strings.TrimPrefix(tokens[i], "INT:")
						//convert to int
						max, _ = strconv.Atoi(tokens[i])
					}
					i++
				}

				randomNumber := rand.Intn(max+1-min) + min

				return strconv.Itoa(randomNumber), i
			} else if tokens[i] == "Read" {
				i++
				i++
				return lib.Read(), i
			}
		} else {
		}

		i++
	}
	return "", i
}

// getbool function to evaluate boolean expressions
func getbool(tokens []string, i int, vars map[string]string, fun map[string][]string) bool {
	// List of operators
	ops := []string{"DEQ", "LESS", "MORE", "NOT"}

	// Create a map for the first set of tokens
	toks1 := []string{}
	i2 := 0

	// Collect tokens until an operator is found
	for !lib.Contains(tokens[i], ops) {
		toks1[i2] = tokens[i]
		i++
		i2++
	}

	op := tokens[i]
	if op == "DEQ" || op == "NOT" {
		i++
	}
	i++

	// Create a map for the second set of tokens
	toks2 := []string{}
	i2 = 0

	// Collect tokens until a newline or opening brace is encountered
	for (tokens[i] != "NEWLINE") && (tokens[i] != "OP_S_B") {
		toks2[i2] = tokens[i]
		i++
		i2++
	}

	// Add "NEWLINE" tokens to both sets
	toks1[len(toks1)] = "NEWLINE"
	toks2[len(toks2)] = "NEWLINE"

	// Get values from the token sets
	val1, _ := getvalue(toks1, 0, vars, fun)
	val2, _ := getvalue(toks2, 0, vars, fun)
	val1i, _ := strconv.Atoi(val1)
	val2i, _ := strconv.Atoi(val2)

	// Compare values based on the operator and return the result
	if (op == "DEQ") && (val1 == val2) {
		return true
	} else if (op == "NOT") && (val1 != val2) {
		return true
	} else if (op == "LESS") && (val1i < val2i) {
		return true
	} else if (op == "MORE") && (val1i > val2i) {
		return true
	}
	return false
}

// interpret function to process and execute the interpreted code
func interpret(tokens []string) {
	// Initialize a map to store functions
	fun := make(map[string][]string)

	// Initialize a map for instance variables
	ivars := make(map[string]string)

	// Initialize the index variable
	i := 0

	// Loop through the tokens
	for i < len(tokens) {
		// Initialize a map for the current function's code
		funcode := []string{}

		// Check if the current token indicates the start of a function
		if tokens[i] == "KEYWORD:func" {
			// Initialize variables for function code extraction
			i2 := 0
			i3 := i

			// Go to an opening brace
			for tokens[i3] != "OP_B" {
				i3++
			}
			i3++
			i4 := 0
			for tokens[i3] != "CL_B" {
				// This process func parameters
				if strings.HasPrefix(tokens[i3], "STRING:") {
					funcode[i4] = "VAR:" + strings.TrimPrefix(tokens[i3], "STRING:")
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
		} else if tokens[i] == "KEYWORD:import" {
			// If the token indicates an import statement
			// Create a new Interpret instance for the imported file
			in := Interpret{
				KeyWords: []string{"print", "if", "var", "func", "while", "import", "return"},
			}

			// Move to the next token until a file path is found
			for !strings.HasPrefix(tokens[i], "TEXT:") {
				i++
			}
			filename := strings.TrimPrefix(tokens[i], "TEXT:")
			// Get the file path and read the content of the file
			executablePath, err := os.Executable()
			if err != nil {
				lib.Print("Nelze získat cestu k spustitelnému souboru:" + err.Error())
				return
			}
			file_path := filepath.Dir(executablePath) + "/" + "libs/" + strings.TrimPrefix(tokens[i], "TEXT:") + ".v"
			data, _ := os.ReadFile(file_path)

			// Replace line endings and tokenize the content
			input := string(data)
			input = strings.Replace(input, "\r\n", "\n", -1)
			in.lexer(input)

			// Loop through the tokens of the imported file
			n := 0
			for n < len(in.tokens) {
				// Check if the token indicates the start of a function in the imported file
				if in.tokens[n] == "KEYWORD:func" {
					// Initialize variables for function code extraction
					i2 := 0
					i3 := n

					for in.tokens[i3] != "OP_B" {
						i3++
					}
					i3++
					i4 := 0
					for in.tokens[i3] != "CL_B" {

						if strings.HasPrefix(in.tokens[i3], "STRING:") {
							funcode = append(funcode, "VAR:"+strings.TrimPrefix(in.tokens[i3], "STRING:"))
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
func (c *code) Code(tokens []string, fun map[string][]string) string {
	// Initialize the index variable
	i := 0
	// Initialize a map for function code
	funcode := []string{}

	// Loop through the tokens to find and store functions
	for i < len(tokens) {
		// Check if the current token indicates the start of a function
		if tokens[i] == "KEYWORD:func" {
			// Initialize variables for function code extraction
			i2 := 0
			i3 := i

			for tokens[i3] != "OP_B" {
				i3++
			}
			i3++
			i4 := 0
			for tokens[i3] != "CL_B" {

				if strings.HasPrefix(tokens[i3], "STRING:") {
					funcode[i4] = "VAR:" + strings.TrimPrefix(tokens[i3], "STRING:")
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
		case tokens[i] == "WHITESPACE":
			// Ignore whitespace tokens

		case tokens[i] == "NEWLINE":
			// Ignore newline tokens

		case strings.HasPrefix(tokens[i], "KEYWORD:"):
			// Check for keyword tokens
			switch {
			case strings.TrimPrefix(tokens[i], "KEYWORD:") == "print":
				// If the keyword is "print," get the value and print it
				val, vl := getvalue(tokens, i, c.vars, fun)
				lib.Print(val)
				i = vl
			case strings.TrimPrefix(tokens[i], "KEYWORD:") == "var":
				// If the keyword is "var," get the variable name and value, then store it
				fname := getname(tokens, i)
				for tokens[i] != "EQ" {
					i++
				}
				i++
				val, vl := getvalue(tokens, i, c.vars, fun)
				c.vars[fname] = val
				i = vl
			case strings.TrimPrefix(tokens[i], "KEYWORD:") == "if":
				// If the keyword is "if," get the code block and execute it if the condition is true
				i++
				toks, ifl := GetCode(tokens, i)
				if getbool(tokens, i, c.vars, fun) {
					c.Code(toks, fun)
				}
				i = ifl
			case strings.TrimPrefix(tokens[i], "KEYWORD:") == "while":
				// If the keyword is "while," get the code block and execute it while the condition is true
				i++
				toks, ifl := GetCode(tokens, i)
				for getbool(tokens, i, c.vars, fun) {
					toks, ifl = GetCode(tokens, i)
					c.Code(toks, fun)
				}
				i = ifl
			case strings.TrimPrefix(tokens[i], "KEYWORD:") == "return":
				// If the keyword is "return," get the value and return it
				i++
				val, _ := getvalue(tokens, i, c.vars, fun)
				return val
			}
		case strings.HasPrefix(tokens[i], "STRING:"):
			// If the token is a string, check if it corresponds to a function
			tokens[i] = strings.TrimPrefix(tokens[i], "STRING:")
			if _, exists := fun[tokens[i]]; exists {
				// Get the function arguments and prepare for function execution
				i3 := i
				for tokens[i3] != "OP_B" {
					i3++
				}
				i3++
				i4 := 0
				fargs := make(map[int]string)
				for tokens[i3] != "CL_B" && tokens[i3] != "NEWLINE" {
					if tokens[i3] != "COM" {
						//i4b := 0
						fargs[i4], _ = getvalue(tokens, i3, c.vars, fun)
						i4++ //= i4b
					}

					i3++
				}
				// Prepare for function execution
				fnum := 0
				fvars1 := make(map[string]string)
				funcp := make(map[string][]string)
				for key, value := range fun {
					funcp[key] = value
				}
				funCopy := make(map[string][]string)
				for k, v := range fun {
					funCopy[k] = []string{}
					for k2, v2 := range v {
						funCopy[k][k2] = v2
					}
				}
				for strings.HasPrefix(funcp[tokens[i]][fnum], "VAR:") {
					// Assign values to function parameters
					fvars1[strings.TrimPrefix(funcp[tokens[i]][fnum], "VAR:")] = fargs[fnum]
					funcp[tokens[i]][fnum] = "WHITESPACE"
					fnum++
				}

				// Create a new code instance for function execution
				c2 := code{
					vars: fvars1,
				}
				fun = funCopy
				// Execute the function code
				c2.Code(funcp[tokens[i]], fun)
				i = i3
			} else if _, exists := c.vars[tokens[i]]; exists {
				// If the token is a variable, replace it with its value
				fname := getname(tokens, i)
				for tokens[i] != "EQ" && i < len(tokens) {
					i++
				}
				i++
				val, vl := getvalue(tokens, i, c.vars, fun)
				c.vars[fname] = val
				i = vl
			}

		default:
			// Handle unknown keywords
			if tokens[i] == "WHITESPACE" {

			} else if tokens[i] == "NEWLINE" {
			} else {
				lib.Print("unknown keyword")
			}
		}
		// Move to the next token
		i++
	}
	return ""
}

func main() {

	i := Interpret{
		tokens:   make([]string, 0),
		KeyWords: []string{"print", "if", "var", "func", "while", "import", "return"},
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "get" {
			//get if args exist
			if len(os.Args) > 2 {
				//get file path
				executablePath, err := os.Executable()
				if err != nil {
					lib.Print("Nelze získat cestu k spustitelnému souboru:" + err.Error())
					return
				}
				file_name, _ := lib.ExtractFileName(os.Args[2])
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

		}

	} else {
		file_path := ""
		//run if args exist
		if len(os.Args) > 2 {
			//get file path
			file_path = os.Args[1]
		} else {
			cur_dir, _ := os.Getwd()
			file_path = cur_dir + "/main.v"
		}

		data, _ := os.ReadFile(file_path)

		input := string(data)
		input = strings.Replace(input, "\\r\\n", "\\n", -1)
		i.lexer(input)
		interpret(i.tokens)
	}
}
