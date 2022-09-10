# PocoCMS Themes: Technical overview

# TODO: Finish!

Left to itself, PocoCMS turns [markdown](glossary.html#markdown) into
pure HTML with no styling. If you add a theme to the front matter,
that all changes:

```
---
theme: ".poco/themes/tufte"
---
# The Tufte theme
My, aren't we elegant
```

## Quick overview: Page themes and global themes 

### Theme

{{- /* Illustration needed */ -}}

Themes contain stylesheet and individual [style tags](glossary.html#style-tags)
for the the page, which consists of the 
[article](glossary.html#article) (body text) and other page
[layout elements](glossary.html#layout-element): the [header](glossary.html#header), 
[nav bar](glossary.html#nav), [aside](glossary.html#aside), and [footer](glossary.html#footer).

The previous example causes the current page to use the factory-installed
Tufte theme. It does not affect themes used by other pages. For that, you'd
use a [global theme](#global-theme).

Because this works only on the current page, it's often called the *page theme*.

### Page theme

*Page theme* is synonomous for themes set using the front matter like this:

```
---
theme: ".poco/themes/tufte"
---
```

It's called that because it only affects that page.

### Global theme

You can use a *global theme* to create a default theme for all
pages on your site. To do so, just use `global-theme` on
your site's [home page](glossary.html#home-page), which is
in the root directory and is named eitheer `README.md` or
`index.md`:


```
---
global-theme: ".poco/themes/pocodocs"
---
Everything's coming up Poco
```

Notes about global themes:

* If you use `global-theme` on any page other than the home page,
it will be ignored
* Declaring a page theme (that is, `theme:` in the front matter)
overrides the global theme for that page--even the home page.

## Example theme file

```
---
author: "Tom Campbell"
branding: "PocoCMS documentation theme"
header: "header.md"
nav: "nav.md"
aside: "aside.md"
footer: "footer.md"
stylesheets:
- "../../css/reset.css"
- "../../css/pococms.css"
- "pocodocs.css"
styletags:
- "article>p{color:gray;}"
---

```


## Contents of a theme

A theme is named by its directory, for examle, 
`.poco/themes/pocodocs`.


### Minimal theme requirements 

A theme directory contains at a minimum

* A `LICENSE` text file ([GitHub](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/licensing-a-repository) has a good overview of common license types)
* A `README.md` file. Its body text should be a paragraph or less describing the theme. Its front matter should point to components that
make up the theme: some combinatino of stylesheets, Markdown files,
HTML files, stylesheets required, and style tags. 
Technically it may be empty. 

### Common theme components

The theme's README.md will normally contain one or all of these components:
* A list of zero or more stylesheets.
