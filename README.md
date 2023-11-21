# Voca
Voca is a simple programming language, written in Go, which I am currently working on. The purpose of the project is to create a language that's easy to understand and write for a beginner (and I'm mainly learning Go on this project).

**TODO:**
 - [x] Creating and calling functions
 - [x] Parameters in functions
 - [x] A function can return a value
 - [x] Variables  
 - [x] Print() function 
 - [x] If and While
 - [x] Import other files written in Voce
 - [x] String concatenation
 - [ ] Several simple libraries (e.g., library for advanced math
       functions, etc.)
 - [ ] Possibility to create a GUI


**How to setup:**
1. Download executeble for your system
2. Write code and save it as main.v (in future there will be command lika run path/to/program.v)
3. Run interpreter.


## Hello World Program:
```go
func main(){
    var s = "Hello World"
    print s
}
```

**Now let's break down the code:**
* All code must be in the main() function.
* Variables do not have types, they are defined using the var keyword.
* We can print the value using print.
