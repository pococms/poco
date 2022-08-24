# Front Matter

*Front matter* optionally starts your file. It contains
instructions for things Markdown can't do for you,
for example, choosing a theme, inserting Javascript
or new style tags into your document, and ensuring
each page can have a unique `<title>` tag.

This page gives you a somewhat technical overview
of front matter, then explains all the 
front matter options PocoCMS provides.

[Formatting rules](#formatting-rules)  
[Front matter basics](#front-matter-basics)   

## Alphabetical
[Author](#author)  
[Description](#description)  
[Keywords](#keywords)  
[Key/value pairs](#keyvalue-pairs)  
[Robots](#robots)  
[Skippublish](#skippublish)  
[Title](#title)  

## Front Matter basics

The simplest PocoCMS Markdown document looks something
like this:

    hello, world.

It generates the following HTML document:

```html
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Powered by PocoCMS</title>
</header>
	<article id="article">
  <p>hello, world.</p>
	</article>
</div><!-- content-wrap -->
	</div><!-- page-container -->
</body>
</html>
```

You will normally see a more complex structure, with
nonprinting commands that start the file:

```yaml
---
Title: "Introducing PocoCMS"
Keywords: "static site generator, jamstack, cms"
---
hello, world. 
```

The text between the two `---` lines is considered a separate
document with no direct relation to the remainder of the
file, which is the Markdown portion. It's called *front matter*
by convention, but what lies between the bracketing `---` 
characters is essentially a database in [YAML](https://yaml.org/)
format.  This data is used to control the format and
output of the HTML files PocoCMS produces.

Here's what happens when PocoCMS generates HTML for the
previous example:

```html
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Introducing PocoCMS</title>
  <meta name="keywords" content="static site generator, jamstack, cms">
</header>
	<article id="article">
  <p>hello, world.</p>
	</article>
</div><!-- content-wrap -->
	</div><!-- page-container -->
</body>
</html>

```

In the first example PocoCMS inserted this shamelessly 
self-promoting `<title>` tag:

```html
<title>Powered by PocoCMS</title>
```

That's because an HTML file is not considered valid without a title. 
As you've already surmised, this portion of the front matter

```yaml
---
Title: "Introducing PocoCMS"
---
```

Is responsible for this line of HTML:

```html
<title>Introducing PocoCMS</title>
```

And obviously this line of front matter

```yaml
---
Keywords: static site generator, jamstack, cms
---
```

Caused a `keywords` metatag to be inserted into the file:

```html
<meta name="keywords" content="static site generator, jamstack, cms">
```

## Formatting rules

* The front matter always starts with a line consisting
solely of 3 dashes: `---`
* If you include front matter in a Markdown file, the
first line must be `---` and cannot be blank.
* The front matter always ends with a line consisting
solely of 3 dashes: `---`
* The dashed lines demarcate YAML front matter but do not include it.
* The key name (on the left side, behind the colon) is
case sensitive. So while this will create a `<title>`
tag in your HTML document:

```yaml
---
Title: PocoCMS makes cry out with joy
---
```

This will not:

```yaml
---
title: I'm totally invisible
---
```

* YAML documents usually can't be empty. Don't start a file
with empty front matter like this:

```yaml
---
---
```

## Author 

Causes an `author` metatag to be inserted into the file.

### Example

Using this `Author` declaration in the front matter:

```yaml
---
Author: "Tom Campbell"
---
```

Causes this metatag to be generated:

```html
<meta name="author" content="Tom Campbell">
```

## Description

Causes a `description` metatag to be inserted into the file.

### Example

Using this `Description` in the front matter:
```yaml
---
Description: "PocoCMS is the easiest static site generator available"
---
```

Causes this metatag to be generated:

```html
<meta name="description" content="PocoCMS is the easiest static site generator available">
```


## Keywords

Causes a `keywords` metatag to be inserted into the file.

### Example

Using these `Keywords` in the front matter:
```yaml
---
Keywords: "static site generator, jamstack, cms"
---
```

Causes this metatag to be generated:

```html
<meta name="keywords" content="static site generator, jamstack, cms">
```


## key/value pairs 

* The front matter contents always start with a key name,
such as `Title` immediately followed by a colon.
* The value associated with it must follow. It's best to use
quotes around the value because stylesheets use characters
that confuse the YAML processing.
* In the YAML below, `Title` is the key, and `Welcome to PocoCMS`
is the value:

```yaml
---
Title: "Welcome to PocoCMS"
---
```

* One enormously powerful feature of YAML is that while every key 
must have exactly one value, that value can consist of more than one item,
a compound structure like a database record, or even an entire database
consisting of multiple records.

In the YAML below, `Stylesheets` is the key. The mutliple
itmes below it are considered a YAML [list](https://docs.ansible.com/ansible/latest/reference_appendices/YAMLSyntax.html#yaml-basics). The list
has one name (`Stylesheets`) but the value  has multiple items in it.
The value in this case is the list consisting of `["poquito.css", "pococms.css"]`.

```yaml
---
Stylesheets: 
- poquito.css
- pococms.css
---
```

This page details all front matter options.

## Robots 

Causes a [`robots` metatag](https://moz.com/learn/seo/robots-meta-directives#:~:text=Robots%20meta%20directives%20(sometimes%20called,or%20index%20web%20page%20content.) to be inserted into the file.

### Example

Using this `Robots` entry in the front matter:

```yaml
---
Robots: "NoIndex"
---
```

Causes this metatag to be generated:

```html
<meta name="robots" content="noindex">
```

**NOTE**

Be careful following this example! It tells 
search engines *not* to index your page, which is
the opposite of what you normally want.


## SkipPublish
`SkipPublish` lists files and directories you don't want to be published.

Remember that if a directory contains the files `index.md`, `installation.md`,
and `avatar.png`, and `401K-info.xls`, here's what will happen when you 
run PocoCMS:

* `index.md` and `installation.md` will be converted to HTML document files
named `index.html` and `installation.html`. The the HTML files will be published
The Markdown files will not.
* By design, all files other than than Markdown files get published. That's
because they're assumed to be required for the site. In this example,
it makes sense that `avatar.png` gets published but you may not want
the personal spreadsheet with your 401K details published.
* Likewise, you probably don't want directories like `node_modules`,
`.DS_Store`, or `.git` published. 
`.git` is included here for good form but 
PocoCMS treats directories with that name starts with `.` as hidden and
doesn't publish them.
* The answer to these problems is to list what you don't want published in `SkipPublish` 
as shown below.

### Example

```yaml
---
SkipPublish:
- 401k-info.xls
- node_modules
- www
- .git
- .DS_Store
- .gitignore
---

```

Here are some common items to skip:

```yaml
---
SkipPublish:
- node_modules
- htdocs
- public_html
- www
- .git
- .DS_Store
- .gitignore
---
```


## Title
The `Title` key lets you set a title for your HTML page.
This has a number of important benefits.

* It will be used for browser tabs open to that page
* It may influece search results
* It assists screen readers for visually impaired users
* It allows you to create a unique title for each page, which
is considered table stakes for government accessibility requirements
in the USA
* It is one of the few tags required to make an HTML-conformant document

Example:

```
---
Title: "Static generator overview"
---

Here's your {{ .Title }}.
```

This page would have its `keywords` [metatag](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name) set to 
`Static generator overview` and the page generated would
read as shown beloe in a web browser:


```
Here's your Static generator overview
```

## Stylesheets 

Causes a `<style>` to be inserted into the file
for each file in the list.

### Example

Using these `Stylesheets` in the front matter:
```yaml
---
Stylesheets: 
- "poquito.css"
- "https://cdn.jsdelivr.net/gh/pococms/poco/pages/assets/css/pocodocs.css"
---
```

Causes this HTML to be generated:

```html
<link rel="stylesheet" href="poquito.css">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/pococms/poco/pages/assets/css/pocodocs.css">
```


