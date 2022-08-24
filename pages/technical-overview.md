---
Keywords: "PocoCMS, PocoCMS summary, technical overview, how does PocoCMS work"
---

# PocoCMS Technical Overview

PocoCMS is a single executable file that reads a
directory tree of files, converts Markdown files
to HTML, and passes the rest through to
the webroot directory unchanged. 
The [webroot](glossary.html#webroot) is where
HTML, stylesheets, and other file assets
are sent to be published on the Web.

PocoCMS is meant above all to be unobtrusive,
easy to learn, and extremely easy to get
started with. You don't
have to create weird special files get started,
or download a theme from some obscure
location on the web. Just type some Markdown, making
sure the root of your site has a file named `index.md`
or `README.md`, and you can get started immediately.

Even if you don't know [Markdown](glossary.html#markdown),
you can just type plain text. Even links get turned into
live hyperlinks without any effort.

## Looking good

PocoCMS has theme support, so you can just mention a theme
in the [front matter](glossary.html#front-matter) of
the [home page](glossary.html#home-page) and
all other pages in the site will inherit that theme.
A theme can also include template files for the
header, nav bar, aside (a.k.a. sidebar), and footer.

## Running PocoCMS

Suppose you plan to create a new site, also known as a project.
You'd create a subdirectory, in this example `mysite`, and
make it the current directory:

```
# Create a subdirectory named mysite. 
# Replace with your own directory name here.
mkdir ~/mysite
# Make it the current directory.
cd ~/mysite
```

* Run poco:

```
poco
```

You're then informed:

```
Site published to /Users/tom/pococms/poco/ed/WWW/index.html
```

Here's what happens when you run `poco` with no 
[command-line options](cli.html) as shown previously:

* PocoCMS looks around to see if the directory is empty (as is 
the case in this example, because it was just created).
* If that directory was empty, and if it has
no [home page](glossary.html#home-page), PocoCMS
generates a simple `index.md` file. If you view its contents,
you'll see this:

```markdown {hl_lines=["1,2]"}
---
Stylesheets:
    - "https://cdn.jsdelivr.net/gh/pococms/poco/pages/assets/css/poquito.css"
---
# Welcome to mysite

hello, world.

Learn more at [PocoCMS tutorials](https://pococms.com/docs/tutorials.html) 

```

{{- /* TODO: Add screenshot */ -}}

When you load the web page you'll see it has minimal styling
and that the first part of the file doesn't get displayed.
That first part between the `---` lines
is called the [front matter](front-matter.html),
and it provides a set of instructions to format and display
the generated HTML, which starts with **Welcome to mysite**.

The front matter is considered a separate document.
It is in [YAML](https://yaml.org) format. The only YAML
in this document is `Stylesheets:`, which is used
to generate a `<link>` tag that pulls in a minimal
spreadsheet from a CDN. It could just as easily
be a local file but this makes for a quick demo.

## PocoCMS themes

