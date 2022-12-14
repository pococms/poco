---
title:  "PocoCMS"
description: "PocoCMS builds fast websites fast"
keywords: "static site generator,ssg,jamstack,cms,ghost.org,gohugo.io"
theme: "pocodocs"
skip:
- "yo mama"
---
# Poco CMS, the world's easiest static site generator

## References
* The PocoCMS [Reference](pages/reference.html) page
* The PocoCMS [FAQ](pages/faq.html)
* Poco CMS [command line options](pages/cli.html)


![](https://www.youtube.com/watch?v=dQw4w9WgXcQ)

*13 2022* 19:00 Many bugs to go. Soft opening on 1 October 2022.

[Report an issue](https://github.com/pococms/poco/issues)

## Building from source

For the moment, you need to build PocoCMS yourself as a 
Go program. Don't worry. There are explicit instructions at
[Build from source](pages/build-from-source.html)

## To create a website using PocoCMS

* Create a destination directory for your project and change to it:

```
mkdir ~/mysite
cd ~/mysite
```

* Create the home page and build the site.

```
poco
```

Upon completion Poco tells you where to find
the generated HTML:

```
Site published to /Users/tom/mysite/WWW/index.html
```

You can open that page in a web browser to get a pretty decent
idea of what it looks like. For a totally functional version
you'll need to see it running on a web server, even `localhost`.


## Creating pages
[Front Matter](pages/front-matter.html)  
[CSS tips](pages/css-tips.html)  



