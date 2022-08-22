---
Title:  PocoCMS

Theme: pages/themes/probot
#Theme: pages/themes/newpoquito
#Theme: pages/themes/tufte
#BAD Theme: pages/themes/poco
#Theme: pages/themes/simplicity
#Theme: pages/themes/latest
Nav: SUPPRESS
Footer: SUPPRESS
#StyleFiles:
#- pages/themes/newpoquito/dark.css
# - poquito.css1
# - pages/assets/css/pococms.css
---

## To build from source:

Currently PocoCMS must be built from source using Go.
It's easy to install Go and even to build, because
PocoCMS is a single file.

    $ # Create a directory. It doesn't have to be here.
    $ mkdir ~/pococms
    # Navigate to that directory.
    $ cd ~/pococms
    $ # Clone the repo.
    $ git clone https://github.com/pococms/poco
    $ # The repo is now in ~/pococms/poco, so navigate there.
    $ cd poco
    $ # And compile: 
    $ go build 
    $ ### OR....
    $ # There's only one file, so you can also use go run.
    $ # That runs the go compiler, then executes the program 
    $ # if there are no compilation errors.
    $ go run main.go
    $ # This will generate an example file
    $ ./poco
    # (Then make sure poco is on your path)

## To create a website using PocoCMS
    $ mkdir ~/mysite
    $ cd ~/mysite
    $ # Create the home page
    $ poco

Poco will create a home page for you. It's `index.md`.


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


