---
Title: README
Header: pages/home/header.md
#Nav: ./pages/diagnostics/layout/nav.md
Aside: ./pages/diagnostics/layout/aside.md
Footer: ./pages/diagnostics/layout/footer.md
LinkTags:
- "<link rel='icon' type='image/png' sizes='32x32' href='/favicon-32x32.png'>"


Setup: [                             
  StyleFileTemplates: [    
    "/pages/assets/css/poquito.css",
    "/pages/assets/css/pococms.css"
  ],
  StyleListTemplates: ["/pages/assets/css/poquito.css", "/pages/assets/css/pococms.css"],
  UnquotedListTemplates: [
    /pages/assets/css/poquito.css, 
    pages/assets/css/pococms.css
  ],
  arr: [ 1, 2, 3, "four" ],
  dict: [ foo: bar]
]


StyleFileTemplates:
  - /pages/assets/css/poquito.css
  - /pages/assets/css/pococms.css

StyleTags:
  - "article{margin-left:12em;margin-right:5em;background-color:white;}"
---

# {{ .Title }}

**Site config**

{{ .Site }}



## Code fences

    ---
    print "4 spaces in"
    ---

     ---
     print "5 spaces in"
     ---


        ---
        print "8 spaces in"
        ---



## Creating pages
[Front Matter](pages/front-matter.html)  
[CSS tips](pages/css-tips.html)  

## Themes

* [Probot theme](pages/themes/probot.html) and [Left sidebar version](pages/themes/probot-left.html)
* [Simplicity](pages/themes/simplicity.html) is 
HTML minimalism at its purest, 
with no header, footer, nav, or aside.

## Tools
* Amazing [Favicon generator](https://realfavicongenerator.net) took a PNG image, then turned it into
a wide variety of [favicons](https://en.wikipedia.org/wiki/Favicon).

## HTML references

* [HTML Color Names](https://htmlcolorcodes.com/color-names) for people who like a simplified color chart and who like using actual names for colors

## Web page performance

* Google's [web.dev](https://web.dev/measure/) page quality measurement tool
* The [Pingdom Website Speed Test](https://tools.pingdom.com/) produces the clearest results

## Diagnostics

* Page showing all [FrontMatter settings](pages/diagnostics/allfeatures.html)


