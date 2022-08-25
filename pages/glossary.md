# Glossary


## code fence
A [code fences](highlighting.html) surrounds arbitrary text with lines
consisting of 3 tickmarks: \`\`\` so that the text displays
in a monospace font. It's good for distinguishing blocks
of code in an article. Here's an example.

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
    Title: "Introducing PocoCMS"
    Stylesheets: 
    - poquito.css
    - pococms.css
    ---

See [Front matter](front-matter.html) for more details.

## home page
The home page is a file named either `index.md` or `README.md`
in the root directory of your project. It has some special
qualities, for example, it's the only file you can use to 
set a theme for the site overall.

*README.md vs index.md*

If you have two home page files, one named `README.md` and
another named `index.md`, the one named `README.md` takes
priority. It is renamed `index.html` in the [webroot](#webroot)
directory when your site is generated.

*Why index.html is important*

When someone visits your website, the web server looks specifically
for the distinguished file named `index.html` in its own webroot directory, so
`index.html` has special importance.

The reason `README.md` takes priority over `index.md` is that's how many
previous site generators roll, such as the one on GitHub.

## Layout element
The structure of a  
[complete HTML document](https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Document_and_website_structure#HTML_layout_elements_in_more_detail) 
is based on these tags: `<header>`, `<nav>`, `<aside>`, `<article>`, and `<footer>`. They are also known as *layout elements*.
PocoCMS takes their corresponding tags from the
[front matter](#front-matter)
and uses those rules to generate the contents of each tag.

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

Technically speaking HTML is also a markup language(https://en.wikipedia.org/wiki/Markup_language) but in the context of static site generators
such as PocoCMS the term normally refers to Markdown.

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
A  PocoCMS site can have an optional theme, which is a collection of stylesheets and Markdown files structured in a particular way. A theme has its own folder, which is used as the name of the theme. The theme can specify stylesheets
to include on every page of the site. The theme can also specify 
a header, nav bar, footer, or aside to include on each page.

## web root
Synonymous with [webroot](#web-root).

## webroot
The webroot is a directory contains all files generated by Poco CMS required 
for your website. By default it's a subdirectory under the directory 
used for your home named `WWW` but you can 
designate a different directory using the `webroot`
[command line option](cli.html#webroot)



