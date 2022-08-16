
# CSS tips

# TODO: write Downloadable Fonts section, StyleTags section

In case you're new to CSS, here's some handy code
you can drop right into your own themes.

## Contents
[Font stacks](#font-stacks)  
[Downloadable fonts](#downloadble-fonts)


<a name="font-stacks"></a>
## Font stacks for sans-serif, serif, and monospace 

HTML defines [font families](https://developer.mozilla.org/en-US/docs/Web/CSS/font-family) such as `sans-serif`, `serif`, and `monospace`.
The standard doesn't go much beyond specifying their rough outlines.

Over the years browser and operating system have accreted their own
[contributions](https://en.wikipedia.org/wiki/Core_fonts_for_the_Web) to help 
make better fonts available to a wide audience.

Since there's no way to know exactly what fonts are on your system, you 
can specifiy a `font-family` that specifies the order in which you want
these fonts to be used if they're available. The great thing about `font-family`
is that you're always guaranteed some kind of useful minimum.

Font families have page size overhead, unlike [downloadable fonts](#downloadable-fonts).

Here's what we like to use. It's all opinion, and this list is 
subject to change.

### Font stack for `sans-serif`

Here's the CSS you'd use to specify better defaults for the most
popular font family on the web,`sans-serif`. You can almost
never go wrong with it.

Note that this specifies using this font stack in `article` tags 
but it could just as easily be replaced with 
`h1`, `aside > h2`, etc.

Standard CSS for a sans-serif font stack:

```
article{font-family:"system-ui","-apple-system","LucidaGrandeUI","HelveticaNeueDeskInterface-Regular","HelveticaNeueDeskInterface-Light","DroidSans","Ubuntu Light","Arial","Roboto-Light","Segoe UI Light","Tahoma","sans-serif";}
```

As a [StyleTag](styletag.html):
 
      ---
      StyleTags:
        - article{font-family:'system-ui','-apple-system','LucidaGrandeUI','HelveticaNeueDeskInterface-Regular','HelveticaNeueDeskInterface-Light','DroidSans','Ubuntu Light','Arial','Roboto-Light','Segoe UI Light','Tahoma','sans-serif';}
      ---

### Font stack for a more elegant `serif`

The [Palatino](https://en.wikipedia.org/wiki/Palatino)-style fonts
are elegant and formal-looking. Many systems have something
approximating Palatino.

Note that this specifies using this font stack in `article` tags 
but it could just as easily be replaced with 
`h1`, `aside > h2`, etc.

Standard CSS for Palatino-style font stack:

```
article{'Palatino Linotype','Palatino','Georgia','Times','Times New Roman','New York','serif';}
```

As a [StyleTag](styletag.html):

      ---
        - article{'Palatino Linotype','Palatino','Georgia','Times','Times New Roman','New York','serif';}
      ---
 
### Font stack for `monospace`

The [monospace](https://en.wikipedia.org/wiki/List_of_monospaced_typefaces)-style fonts are unique in that all characters in the font are the same width.
They're good for code listings like the one on this page, for example,
normally found in the HTML `<pre>` and `<code>` elements.

Note that this specifies using this font stack in `article` tags 
but it could just as easily be replaced with 
`h1`, `aside > h2`, etc.

Standard CSS for monospace font stack:

```
article{font-family:Consolas,Monaco,Menlo,'DejaVue Sans Mono','Lucida Console',monospace;}
```

As a [StyleTag](styletag.html):

      ---
        - article{font-family:Consolas,Monaco,Menlo,'DejaVue Sans Mono','Lucida Console',monospace;}
      ---
 
## Downloadable Fonts 

### TODO: Write this section ;)

