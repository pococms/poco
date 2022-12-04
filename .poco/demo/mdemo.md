# PocoCMS Markdown quick reference

This page serves as a Markdown quick reference.
It is also gives you a good idea of how a
classless CSS framework can operate in the
real world.

## Table of contents 

* [Common text formatting](#common-text-formatting)
* [Links](#links)
* [Bookmarks](#bookmarks)
  - [Linking inside a document](#linking-inside-a-document)
  - [All headers are automatically bookmarks](#automatic-bookmarks)
  - [Bookmarks must be unique in an HTML document](#bookmarks-unique)
  - [How to create bookmarks manually](#how-to-create-bookmarks-manually)
  - [Linking to bookmarks on other webstes](#linking-other-websites)
* [Header styles](#header-styles)
* [Coding styles](#coding-styles)
  - [Choosing the programming language](#choose-pl)
* [Ordered lists](#ordered-lists-37)
* [Unordered, or bullet lists](#unordered-lists)
* [The "third" list type: definition lists](#def-lists)
* [Creating clickable image links in Markdown](#clickable-images)
* [Tables](#tables)
* [Block quote](#block-quote)
* [Footer](#footer)
* [HTML Forms](#html-forms)



## Markdown syntax

Here's how markdown appears in the **{{.theme }}** theme
## Common text formatting

#### You type:
```
Normal body text, **strong**, ~~strikethrough~~, and with *emphasis*.
```

#### It shows as:
Normal body text, **strong**, ~~strikethrough~~, and with *emphasis*.

Horizontal rule:

#### You type:
```
---
```

#### It shows as:
---

## Links

#### You type:
```
[HTML color names](https://htmlcolorcodes.com/color-names/)
```

#### It shows as:
[HTML color names](https://htmlcolorcodes.com/color-names/)

<a id="handmade-anchor"></a>
## Bookmarks


Suppose you want to link to a particular location inside a page. As long as there's an `id` attribute in the document's HTML you can cause a link to jump directly to that part of the page by specifing the link following a `#` character. 

### Linking inside a document

We'll loosely call them anchors or bookmarks, 
although the HTML specification simply calls it the [id attribute](https://html.spec.whatwg.org/multipage/dom.html#the-id-attribute).

Here's a demonstration. If you type:

```
Jump to the [tables](#tables) section.
```

The result will be this (click the link, then use your browser's Back button to return here):  

Jump to the [tables](#tables) section.

<a name="automatic-bookmarks"></a>

### All headers are automatically bookmarks too

Metabuzz automatically generates an `id` attribute for each header from h1 to h6 by taking the text of the link itself, reducing it to lowercase, and either replacing spaces and other non-letter characters with hyphens, or removing them altogether. If you look at the HTML for this page you'll see the `Tables` header looks like this:

```
<h2 id="tables">Tables</h2>
```

The `Coding styles` header uses a hyphen to replace the space:

```
<h2 id="coding-styles">Coding styles</h2>
```

And the more complicated example of the header named "The 'third' list type: definition lists":

```
<h3 id="the-third-list-type-definition-lists">The "third" list type: definition lists</h3>
```
<a name="bookmarks-unique"></a>
### Bookmarks must be unique in an HTML document

The `id` attribute must be unique within a document. Notice how on this page there are many headers simply called `You type:`? Metabuzz keeps track of them and turns each of them into unique IDs by naming them `you-type-1`, `you-type-2`, and so forth.

### How to create bookmarks manually

Suppose you want a bookmark that's not a header? You can insert one anywhere by starting a Markdown line with the pure HTML code for anchors. (HTML is [allowed in Markdown with a few restrictions](https://spec.commonmark.org/0.29/#html-blocks)).


You type:

```
<a id="jump-here"></a>
```

Then you create a link to it by adding the `#jump-here` portion to a link, which is noted by the web browser but not displayed:

```
[Learn about blockquotes](#jump-here)
```

Try it now: [Learn about blockquotes](#jump-here)

<a id=linking-other-websites></a>

### Linking to bookmarks on other websites

You can also link to an anchor to other websites, if they have anchors. Here's a link to the history of futbol on Wikipedia:

#### You type:
```
[History of futbol](https://en.wikipedia.org/wiki/Association_football#History)
```

#### It takes you right there:

[History of futbol](https://en.wikipedia.org/wiki/Association_football#History)


## Header styles

#### You type:
```
# h1
## h2
### h3
#### h4
##### h5
###### h6
```

#### It shows as:

# h1
## h2
### h3
#### h4
##### h5
###### h6

## Coding styles

You can format text inline as `code` by surrounding text with `` ` ``tick marks`` ` ``, or go block style by enclosing the lines of code in a "fenced code block", which begin and end with 3 tickmarks:

```
You can format text inline as `code`, or go block style:
```
<a name="choose-pl"></a>
### Choosing the programming language

You can specify a color scheme for a particular programming language by including its name after the first 3 tick marks of the code block.

#### You type:
    ```python
    print ("This is a code block")
    ```

#### It shows as:
```python
print ("This is a code block")
```

#### If you're a Go programmer, you type:

    ```go
    fmt.Println("This is a code block")
    ```

#### It shows as:

```go
fmt.Println("This is a code block")
```

## There are 2 or 3 kinds of list types

#### You type:
```
### Ordered lists

1. Ordered lists have numeric sequences
1. Even though you write `1` in Markdown,
1. The numbers display properly on output
```


#### It shows as:
### Ordered lists

1. Ordered lists have numeric sequences
1. Even though you write `1` in Markdown,
1. The numbers display properly on output



### Unordered lists

Unordered lists are normally represented as bullets They're always
the same, whereas ordered lists show as automatically generated sequences of letters or number

#### You type:
```
Reasons people hate bullet lists

* They were traumatized by bad PowerPoint
* Some peple actually like bullet lists
  + You can indent bullet lists
    - Just use tab, then one of the characters `*`, `+`, `-`
  + The `+` isn't required. It's just for clarity
    - Most Metabuzz themes go up to 3 visible levels
    - Any more levels than 3 makes it hard for the reader
      + Therefore the Metabuzz theme framework seldom covers indentation levels as deep as this bullet point
```

#### It shows as:

Reasons people hate bullet lists

* They were traumatized by bad PowerPoint
* Some peple actually like bullet lists
  + You can indent bullet lists
    - Just use tab, then one of the characters `*`, `+`, `-`
  + The `+` isn't required. It's just for clarity
    - Most Metabuzz themes go up to 3 visible levels
    - Any more levels than 3 makes it hard for the reader
      + Therefore the Metabuzz theme framework seldom covers indentation levels as deep as this bullet point

<a name="def-lists"></a>
### The "third" list type: definition lists

A definition list lets you display things like an item
and its meaning in a distinct way:

#### You type:
```
Definition list
: A way to show a visual relationship between a word or phrase
and an explanation for it

Markdown
: A convention for generating HTML from a more human-readable 
source format.
```

#### It shows as:
Definition list
: A way to show a visual relationship between a word or phrase
and an explanation for it

Markdown
: A convention for generating HTML from a more human-readable 
source format.

<a name="clickable-images"></a>
### Creating clickable image links in Markdown

Remember that a Markdown link looks like this:

```
[Twitter](https://twitter.com)
```

And that an image link looks like this:

```
![Twitter logo](twitter-32x32-black.png)
```

You can combine them to make a clickable image, like this:

```
[![Twitter logo](twitter-32x32-black.png)](https://twitter.com)
```

## Tables

Use this method of creating tables. Columns are normally left-aligned,
but `:|` on the row of dashes right-aligns a column, and  `|:-` and  `-:|` center-aligns a column.
Headers are always centered.

#### You type:
```

| Left-justified Contents  |  Centered Contents   | Right-justified Contents   |
| ------------------------ |:--------------------:| --------------------------:|
| Row 1, Col 1             | Row 1, Col 2         | Row 1, Col 3               |
| Row 2, Col 1             | Row 2, Col 2         | Row 2, Col 3               |

```

And here's what results from the table markdown shown above:
#### It shows as:

|  Left-justified Contents |  Centered Contents   | Right-justified Contents   |
| ------------------------ |:--------------------:| --------------------------:|
| Row 1, Col 1             | Row 1, Col 2         | Row 1, Col 3               |
| Row 2, Col 1             | Row 2, Col 2         | Row 2, Col 3               |

## Block quote

#### You type:
```
>Hypocrisy waits silently for us all. 
```

#### It shows as:
> Hypocrisy waits silently for us all.

[Return to the bookmarks section](#bookmarks)

## HTML forms
**HTML form test: Search reddit.com**
<form action="https://www.google.com/search" class="searchform" method="get" name="searchform" target="_blank">
<input name="sitesearch" type="hidden" value="reddit.com">
<input autocomplete="on" name="q" placeholder="Search reddit.com" required="required"  type="text">
<button class="button" type="submit">Search</button>
</form>


## Footer

Return to the [Bookmarks](#handmade-anchor)

