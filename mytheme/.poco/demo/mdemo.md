
# PocoCMS Markdown quick reference

This page serves as a Markdown quick reference.
It is also gives you a good idea of how a
classless CSS framework can operate in the
real world.

## Table of contents 

* [Common text formatting](#common-text-formatting)
  - [Horizontal rule](#horizontal-rule)
* [Links](#links)
* [Bookmarks](#bookmarks)
  - [Linking inside a document](#linking-inside-a-document)
  - [All headers are automatically bookmarks](#all-headers-are-automatically-bookmarks-too)
  - [Bookmarks must be unique in an HTML document](bookmarks-unique#)
  - [Linking to bookmarks on other webstes](#linking-other-websites)
* [Header styles](#header-styles)
* [Coding styles](#coding-styles)
  - [Choosing the programming language](#choosing-the-programming-language)
* [Ordered lists](#ordered-lists-37)
* [Unordered, or bullet lists](#unordered-lists)
* [The "third" list type: definition lists](#the-third-list-type-definition-lists)
* [Creating clickable image links in Markdown](#creating-clickable-image-links-in-markdown)
* [Tables](#tables)
* [Block quote](#block-quote)
* [HTML Forms](#html-forms)


The general principle with Markdown is that things
you might type if you didn't have formatting get
replaced with their formatted HTML equivalent.
What follows is most of the Markdown supported by PocoCMS.

## Common text formatting

#### You type:

```
Normal body text, **strong**, ~~strikethrough~~, and with *emphasis*.
```

#### It shows as:

Normal body text, **strong**, ~~strikethrough~~, and with *emphasis*.

## Horizontal rule

You can get a line (horizontal rule) that stretches across the
page, depending on how it's styled. Just use a line consisting
of nothing but `--`.

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
### Bookmarks must be unique in an HTML document

The `id` attribute must be unique within a document. Notice how on this page there are many headers simply called `You type:`? Metabuzz keeps track of them and turns each of them into unique IDs by naming them `you-type-1`, `you-type-2`, and so forth.


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

Return to the [Bookmarks](#handmade-anchor)

