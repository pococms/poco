---
stylesheets:
- rawdog.css
---

# Rawdog theme for PocoCMS

Rawdog uses the most minimal of style sheets.
It employs plain HTML and just a few tweaks to default
colors, margins, and padding. It's a good them for undistracted writing.
Here's the whole thing.

It does NOT support
* Header
* Nav
* Aside
* Footer

Rawdog's theme consists of this single file, called `rawdog.css`:

```
body{margin:5rem}
article{background-color:white;padding:4rem;border:1px solid gray;box-shadow:3px 3px gray;}
pre{background-color: ghostwhite;font-size: smaller;overflow: auto;padding: .5em}
@media (max-width:720px){
  body{margin:1rem;}
  article{padding:1rem;margin:1rem;}
  html{font-size:20px;}
}
```


