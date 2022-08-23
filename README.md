---
Title:  PocoCMS
Theme: pages/themes/pocodocs
StyleTags:
  - "article{margin:0 8em 0 8em}"
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

## To build from source:

Until 1 September 2022 PocoCMS must be built from source using Go.
It's easy to install Go and even to build, because
PocoCMS is a single file.


### One time: install Go and git

* Install the [Go language](https://go.dev/dl/) if necessary.
* Install [git](https://git-scm.com/downloads) if necessary.
* Create a directory or change to a directory to install PocoCMS.

```
mkdir ~/pococms
```

* Navigate to that directory.

```
cd ~/pococms
```

* Clone the PocoCMS repo

```
git clone https://github.com/pococms/poco
```

* The repo is now in ~/pococms/poco, (in this example) so navigate there.

```
cd poco
```

### One time: compile PocoCMS

* And compile: 

```
go build 
```

**OR...**


* There's only one file, so you can also use go run.
That runs the go compiler on the single `main.go` 
containing PocoCMS, then executes PocoCMS.

```
go run main.go
```


Add the `poco` executable to your system path 
so you can run it from any directory.

### PocoCMS creates a starting project automatically.

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


