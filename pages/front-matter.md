# Front Matter

# TODO: SkipPublish isn't implemented yet

## Front Matter basics

The simplest PocoCMS Markdown document looks something
like this:

    hello, world.

It generates the following HTML document:
```
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

    ---
    Title: "Introducing PocoCMS"
    Keywords: "static site generator, jamstack, cms"
    ---
    hello, world. 

The text between the two `---` lines is considerd a separate
document with no direct relation to the remainder of the
file, which is the Markdown portion. It's called *front matter*
by convention, but what lies between the bracketing `---` 
characters is essentially a database in [YAML](https://yaml.org/)
format.  This data is used to control the format and
output of the HTML files PocoCMS produces.

Here's what happens when PocoCMS generates HTML for the
previous example:

```
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

```
<title>Powered by PocoCMS</title>
```

That's because an HTML file is not considered valid without a title. 
As you've already surmised, this portion of the front matter

    ---
    Title: "Introducing PocoCMS"
    ---

Is responsible for this line of HTML:

```
<title>Introducing PocoCMS</title>
```

And obviously this line of front matter

    ---
    Keywords: static site generator, jamstack, cms
    ---

Caused this metatag to be inserted into the file:

```
<meta name="keywords" content="static site generator, jamstack, cms">
```

## Formating of front matter 

* The front matter always starts with a line consisting
solely of 3 dashes: `---`
* If you include front matter in a Markdown file, the
first line must be `---` and cannot be blank.
* The front matter always ends with a line consisting
solely of 3 dashes: `---`
* The lines demarcate YAML front matter but do not include it.
* YAML documents usually can't be empty. Don't start a file
with empty front matter like this:

```
    ---
    ---
```

### key/value pairs in front matter

* The front matter contents always start with a key name,
such as `Title` immediately followed by a colon.
* The value associated with it must follow. It's best to use
quotes around the value because stylesheets use characters
that confuse the YAML processing.
* In the YAML below, `Title` is the key, and `Welcome to PocoCMS`
is the value:

```
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

    ---
    Stylesheets: 
    - poquito.css
    - pococms.css
    ---

This page details all front matter options.

SkipPublish
: `SkipPublish` lists files and directories you don't want to be published.

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

    ---
    SkipPublish:
    - 401k-info.xls
    - node_modules
    - www
    - .git
    - .DS_Store
    - .gitignore
    ---

Here are some common items to skip:

```
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

key, value
: This a definition of how terms are used in front matter. 
The term `key` is 
the technical name for the first part of key/value pair,
and all front matter entries are key/value.


```
---
Title: "Introducing PocoCMS"
Keywords: "static site generator, jamstack, cms"
---
```

