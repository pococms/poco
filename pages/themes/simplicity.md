---
StyleTags:
  - "body{margin:5em}"
  - "html{background-color:gray}"
  - "article{background-color:white;padding:4em}"
  - "pre{background-color: ghostwhite;font-size: smaller;overflow: auto;padding: .5em}"
---
# Simplicity theme

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

