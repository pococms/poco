

# Highlighting code listings

Markdown has a feature allowing you to display
abitrary amounts of text with no formatting, using
monospaced text. 

Here's a typical example.

```
// Return the current time as a string
func theTime() string {
  t := time.Now()
  s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
  return s
}
```

It's called a [code fence](glossary.html#code-fence). 
PocoCMS supports extensions to that
minimal standard, including support for keyword highlighting
in many programming languages and marking specific lines numbers
in the code for clarity.


[How to create code fences](#code-fences)    
[Programming language support](#programming-language-support)  
[Highlighting lines](#highlighting-lines)   




## Code fences

Do this by using one of these techniques:

1. Enclosing the text between lines with 3 \`\`\` tickmarks
or by indenting the text four spaces.
2. Indenting the text exactly 4 spaces per line.

### Code fences using backticks

Here's an example using technique #1, which is 
surrounding the code with lines containing 3 backticks:

```
    ```
    // Return the current time as a string
    func theTime() string {
      t := time.Now()
      s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
      return s
    }
    ```
```

### Code fences using indentation

Here's an example using technique #2, which is indenting 4 spaces:

        // Return the current time as a string
        func theTime() string {
          t := time.Now()
          s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
          return s
        }



## Programming language support

You can get syntax highlighting for most programming languages
by following the backticks with the name of the language.
Here's an example for Go:

```
    ```go
    // Return the current time as a string
    func theTime() string {
      t := time.Now()
      s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
      return s
    }
    ```
```

The result:

```go
// Return the current time as a string
func theTime() string {
  t := time.Now()
  s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
  return s
}
```

### List of supported languages

You can get a list of supported programming languages
for highlighting [here](https://github.com/alecthomas/chroma#supported-languages), thanks to the tireless work of [Alec Thomas](https://github.com/alecthomas).

```
    ```go {hl_lines=[1,"3-4"]}
    // Return the current time as a string
    func theTime() string {
      t := time.Now()
      s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
      return s
    }
    ```
```

## Highlighting lines

You can also show specified lines highlighted inside a code fence.
Following the language designation you'll use an `hl_lines`
directive.

### How to specify a list of lines to highlight

In the following example, the code fence isn't for 
`go` or `java` or `python` but merely `text`. Here's how
to add the `hl_lines` directive:

* Follow the language designation  (`text` in this example), 
with a space.
* After the space, place this code: `{hl_lines=[]}`
* Now inside the square brackets, insert the lines or 
ranges of lines you want to highlight. For example:
* When you enter a number or a range, use double quotes `"` around it.
* To highlight only line 1: `{hl_lines=["1"]}`
* To highlight a range from lines 7-9 inclusive: `{hl_lines=["7-9"]}`
* To highlight lines 3, then 7 to 9 inclusive: `{hl_lines=["3","7-9"]}` 

Here's a typical example.

```
    ```text {hl_lines=["3","7-9"]}
    Sonnet 116

    Love is not love
    Which alters when it alteration finds,
    Or bends with the remover to remove.
    O no, it is an ever-fixed mark
    That looks on tempests and is never shaken;
    It is the star to every wand'ring bark,
    Whose worth's unknown, although his height be taken."

    ― William Shakespeare
    ```
```

It renders like this:

```text {hl_lines=["3","7-9"]}
Sonnet 116

Love is not love
Which alters when it alteration finds,
Or bends with the remover to remove.
O no, it is an ever-fixed mark
That looks on tempests and is never shaken;
It is the star to every wand'ring bark,
Whose worth's unknown, although his height be taken."

― William Shakespeare
```

### Adding line numbers

You can add line numbers by adding
`{linenos="true"}` after the language
designation, which in the following example
is for Go:

```
    ```go {linenos="true"} 
    // Return the current time as a string
    func theTime() string {
      t := time.Now()
      s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
      return s
    }
    ```
```

Resulting in:


```go {linenos="true"}
// Return the current time as a string
func theTime() string {
  t := time.Now()
  s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
  return s
}
```

### Changing the starting line number

You can alter the starting line number
by adding a `linenostart=10` directive,
where you'd replace `10` with whatever number you like.

```
    ```go {linenos="true",linenostart=20} 
    // Return the current time as a string
    func theTime() string {
      t := time.Now()
      s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
      return s
    }
    ```
```

Resulting in:


```go {linenos="true",linenostart=20}
// Return the current time as a string
func theTime() string {
  t := time.Now()
  s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
  return s
}
```



