---
Title:  PocoCMS
Theme: pages/themes/pocodocs
#Theme: pages/themes/latest
#Aside: aside.md
#Theme: pages/themes/test
#Theme: pages/themes/newpoquito
#Theme: pages/themes/probot
#Theme: pages/themes/tufte
#BAD Theme: pages/themes/poco
#Theme: pages/themes/simplicity
#Nav: SUPPRESS
#Footer: SUPPRESS
#Stylesheets:
#- pages/assets/css/poquito.css
#- pages/assets/css/pococms.css
#- pages/themes/newpoquito/dark.css
# - poquito.css1
# - pages/assets/css/pococms.css
---

Paragraph.  

Paragraph.

Not another.

* visit https://pococms.com
* visit http://pococms.com
* visit www.pococms.com

# PocoCMS: the fastest way to turn Markdown into documentation

PocoCMS is a command-line JamStack utility to let you build
documentation websites instantly. It has anywhere from zero
to a short learning curve, and will start you off by
building a website for you the first time you run it.


PocoCMS is meant to give you the same feeling you had if
you created websites in the early days of the Web, because
it requires very little tooling and gives you strict control
over the output. It speeds things up by using Markdown for
text and lets you build or use themes easily.

If you don't want strict control over the output and just want
to get some documentation onto the Web, PocoCMS is even more
suited to you. Just start writing files in Markup, and PocoCMS
will instantly produce something attractive and functional.

## Building from source

For the moment, you need to build PocoCMS yourself as a 
Go program. Don't worry. There are explicit instructions at
[Build from source](build-from-source.html)

## To create a website using PocoCMS

* Create a destination directory for your project and change it it:

```
mkdir ~/mysite
cd ~/mysite
```

* Create the home page and build the site.

```
poco
```

Poco will create a home page for you and tell you where it is:

```
Site published to /Users/tom/pococms/poco/WWW/index.html
```

You can open that site in a web browser to get a pretty decent
idea of what it looks like. For a totally functional version
you'll need to see it running on a web server, even `localhost`.


## Diagnostics
* Page showing all [FrontMatter settings](pages/diagnostics/allfeatures.html)
* [mdemo.html](pages/demo/mdemo.html) Shows most Markup capabilities
* [CSS validator](https://jigsaw.w3.org/css-validator/#validate_by_input)

## Creating pages
[Front Matter](pages/front-matter.html)  
[CSS tips](pages/css-tips.html)  

## Tools
* Amazing [Favicon generator](https://realfavicongenerator.net) took a PNG image, then turned it into
a wide variety of [favicons](https://en.wikipedia.org/wiki/Favicon).

## HTML references

* [HTML Color Names](https://htmlcolorcodes.com/color-names) for people who like a simplified color chart with actual names for colors

## Web page performance

* Google's [web.dev](https://web.dev/measure/) page quality measurement tool
* The [Pingdom Website Speed Test](https://tools.pingdom.com/) produces the clearest results


