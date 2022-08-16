---
StyleTags:
  - "body{margin:5em}"
  - "html{background-color:gray}"
  - "article{background-color:white;padding:4em}"
  - "pre{background-color: ghostwhite;font-size: smaller;overflow: auto;padding: .5em}"
  - "@media (max-width:960px){header,nav,aside,footer{display:none} article{padding:1em;}body{margin:1em}"
---
# Simplicity theme for PocoCMS

Simplicity uses no external style sheets.
It employs plain HTML and just a few tweaks to default
colors, margins, and padding. Here's the whole thing,
which appears at the top of this page in the front matter:

      ---
      StyleTags:
        - "body{margin:5em}"
        - "html{background-color:gray}"
        - "article{background-color:white;padding:4em}"
        - "pre{background-color: ghostwhite;font-size: smaller;overflow: auto;padding: .5em}"
        - "@media (max-width:960px){header,nav,aside,footer{display:none} article{padding:1em;}body{margin:1em}"
       ---

## Code examples

### Go:

```go
fmt.Println("hello, world.")
```
### C
```c
puts("hello, world.")
```

### (Default)
```
print "hello, world."
```

