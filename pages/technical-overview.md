# PocoCMS Technical Overview


**Come back 1 September 2022 for something more substantial**

## TODO: Page, and entry here, on how to configure Markdown extensions

PocoCMS is a single executable file that reads a
directory tree of files, converts Markdown files
to HTML, and passes the rest through. (The file's
extension determines whether a file is treated as Markdown.)

It's meant above all to be unobtrusive. You don't
have to create some weird special files to get your
website going, or download a theme from some obscure
location on the web. Just type some Markdown, making
sure the root of your site has a file named `index.md`
or `README.md`, and you can get started immediately.
Even if you don't know [Markdown](glossary.html#markdown),
you can just type plain text. Even links get turned into
live hyperlinks without any effort.

## Looking good

PocoCMS has theme support, so you can just mention a theme
in the [front matter](glossary.html#front-matter) and
all other pages in the site will inherit that theme.
A theme also 

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

Here's what just happened:

* PocoCMS looked around and found an empty directory.
* Because that directory was empty, and because it had
no [home page](glossary.html#home-page), PocoCMS
generated a stub `index.md` file. If you view its contents,
you'll see this:





## PocoCMS 
## PocoCMS themes

