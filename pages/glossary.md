# Glossary


## article

The *article* is the main body text of a web page. In the
example below, the article consists of the worlds `hello, world.`:

```
---
theme: ".poco/themes/pocodocs"
---
hello, world.
```

When rendered as an HTML page using the PocoDocs theme, you'll also
see a header, nav, aside, and footer. Those are *not* part of
the article. They're known as page [layout elements](#layout-element).

## code block

Synonomous with [code fence](#code-fence)

## code fence

A [code fence](highlighting.html) surrounds arbitrary text with lines
consisting of 3 tickmarks: \`\`\` so that the text displays
in a monospace font. It's good for distinguishing blocks
of code in an article. Here's an example.

```
    // Return the current time as a string
    func theTime() string {
      t := time.Now()
      s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
      return s
    }
```

## CommmonMark
The term *CommonMark* is the name of a community standard for
for the [Markdown](#markdown) text formatting
conventions used to generate your web pages. 
In these help pages it is synonomous with 
Markdown and markup.
<a name="front-matter"></a>

## front matter
Front matter is everything in between the line starting
with `---` to its matching `---`  line at the end. All front
matter entries consist of a single key followed by a `:`
colon character, a space, and a value consisting
of one or more items. PocoCMS uses the key to look up 
the value.

In the example
below, there are two entries in the front matter: `Title`
with the value `Introducing PocoCMS`, and 
has a key nameed `Stylesheets`, with the value consisting
of multiple stylesheet names that will  be added as
separate `<style>` tags in the finished HTML document. 

    ---
    title: "Introducing PocoCMS"
    theme: ".poco/themes/probot"
    stylesheets: 
    - poquito.css
    - pococms.css
    ---

See [Front matter](front-matter.html) for more details.

## global theme

A global [theme](#theme) creates 
default styling for every page
in your site. 

See also [page theme](#page-theme)


## home page

The home page is a file named either `index.md` or `README.md`
in the root directory of your project. It has some special
properties, for example, it's the only file you can use to 
`globaltheme` to set a theme for the site overall.

*README.md vs index.md*

If you have two home page files in the root directory, 
one named `README.md` and another named `index.md`, 
the one named `README.md` takes priority. 
It is renamed `index.html` in the [webroot](#webroot)
directory when your site is generated.

*Why index.html is important*

When someone visits your website, the web server looks specifically
for the distinguished file named `index.html` in its own webroot directory, so
`index.html` has special importance.

The reason `README.md` takes priority over `index.md` is that's how many
previous site generators roll, such as the one on GitHub.

### Defining a global theme on the home page

The home page lets you define a global [theme](#theme) for the entire site.
If you add `globaltheme:` followed by the theme name to the front
matter as shown below, all pages of your site will default to
the global theme without your having to specify it on each page.
For example, this defines `wide` as the global theme:

    ---
    theme: ".poco/themes/wide"
    ---


## Layout element

A finished PocoCMS web page includes the following
layout elements: [header](#header), [nav](#nav),
[article](#article), [aside](#aside), and [footer](#footer). 
Each layout element directly corresponds
to an HTML tag. Most of them can be disabled on a per-page
basis, overriding the theme definition.

### header

The `<header>` element, normally referred to simply as the *header*,
appears at the top of the page. It is likely to look similar on
most pages of your site. It usually makes your site easily
identifiable, normally has a clickable logo that brings
users back to the home page, and may have some common navigation
elements.

### nav

The `<nav>` element, normally called the *nav* or *navbar*, 
is sandwiched between the header and the article. It should
look similar on most pages of your site. It usually has
some common navigation elements.

### article
* The `<article>` element contains the text of your 
Markdown page after conversion to HTML. It appears under the navbar. It is normally unique on each page of your site. Search engines
don't like to see articles or [title tags](front-matter#title)
repeated.

### aside 
* The `<aside>` element acts as a sidebar. It normally appears to the
left or right of the article. HTML recognizes only one aside
per page.

### footer

* The `<footer>` appears at the bottom of the page right after
the article. It should look similar on most pages of your site.
It usually has identifying information for the company and
some common navigation elements such as links to contact,
terms and conditions, privacy policy, and sitemap.

Layout elements appear on the page only if the theme has defined
theme. Most of them can be omitted one page at a time
by using "SUPPRESS" as their value in the front matter

For example, to prevent a sidebar from appearing on the
current page, you'd add this to the [front matter](#front-matter).

     ---
     Aside: "SUPPRESS"
     ---

For more on layout elements, read about the structure of a  
[complete HTML document](https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Document_and_website_structure#HTML_layout_elements_in_more_detail) on MDN.


## Markdown
Markdown is a sensible way to represent text files so that they read easily as plain text if printed out as is, but which also carry enough semantic meaning that they can be converted into HTML. Markdown is technically known as a *markup langauge*, which means that it contains both text, e.g. hello, world, and easily distinguishable annotations about how the text is used, e.g. marking up *hello* to emphasize the word in italics--its markup. The name markdown is a play on the term markup. The name markdown is a play on the term markup. 

The closest thing to an industry standard for Markdown is CommonMark. PocoCMS converts all CommonMark text according to specification, and includes extensions for things like tables, strikethrough, and autolinks. See the source to Goldmark for more information on extensions.

Take this example of Markdown you might use in a document:
```
# Introduction

*hello*, world
```
The above would be converted in HTML that looks like this.
```
<h1>Introduction</h1>
<p<em>hello</em>, world.</p>
```
That means the `# Introduction` line actually represents the HTML heading type h1, which is the hightest level of organization. `## Introduction` would generate an h2 header, and so on.

The asterisk characters are replaced by the `<em>` tag pair, which means they have the semantic power of emphasis. This is represented by HTML as italics, although you could override it in CSS.

In these help pages Markdown is synonyomous with markup and CommonMark.


## Markup 
The term *markup* generally refers to the [Markdown](#markdown) text formatting
conventions used to generate your web pages. In these help pages it is synonomous with 
Markdown, markup, and [CommonMark](#commonmark).

Technically speaking HTML is also a [markup language](https://en.wikipedia.org/wiki/Markup_language) but in the context of static site generators
such as PocoCMS the term normally refers to Markdown.


## page theme

A page [theme](#theme) controls the appearance of a single page. 
It overrides the [global theme](#globaltheme), if any.



## project
A PocoCMS *project* is a directory tree with th
source Markdown files and other assets required to 
create a website. At a minimum, it needs a
[home page](#home-page), which is a Markdown file named
either `index.md` or `README.md`, and a webroot subdirectory,
by default named `WWW`.


## site
See [project](#project)


## source file
A *source file* is the Markdown file used to create a matching HTML file for output.
For example, most directories have a source filenamed `index.html`, which
is the default location web servers look when users navigate to a 
website

## theme

A PocoCMS theme is a collection of stylesheets and Markdown files 
contained in a directory (folder). 
The directory is used as the name of the theme. 
The theme can specify styles to include on every page of the site. 
The theme can also specify 
[layout elements](#layout-elements): a [header](#header), [nav bar](#nav), [footer](#footer), or [aside](#aside) to include on each page.

There two kinds of themes: global, and page. A global theme causes all
pages in the site to use the same theme without having to specify it
every time in the page front matter. See [home page](#home-page) for
a usage example.

You can specify themes on a per page basis. For example, if you want
to use the theme named `wide` you would add this to your Markdown page:

    ---
    theme: .poco/themes/wide
    ---

### How to find out what themes are installed

To find out what themes are installed on your machine, just run
this at the command line:

```bash
poco -themes
```

## web root
Synonymous with [webroot](#web-root).

## webroot
The webroot is a directory contains all files generated by Poco CMS required 
for your website. By default it's a subdirectory under the directory 
used for your home named `WWW` but you can 
designate a different directory using the `webroot`
[command line option](cli.html#webroot)



