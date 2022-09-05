# Using PocoCMS themes

The easiest way to style the appearance of your site
is to use themes. PocoCMS has two kinds of themes: global 
and page.

A [global theme](glossary.html#global-theme) creates 
default styling for every page
in your site. A [page theme](glossary.html#page-theme) applies only to 
that page, so if it's named in the front matter of your page, 
the global theme, if any, is ignored on that page.

## Global theme

A global [theme](glossary.htmle#theme) creates 
default styling for every page
in your site. It's specified using `global-theme`
in the front matter. It can only be used on your
[home page](glossary.html#home-page). 

To apply a global theme named "pocodocs" to
all pages on your site by default:

```yaml
    ---
    global-theme: "pocodocs"
    ---
```


See also [page theme](#page-theme)


## Special rule

A page theme overrides the global theme. 
The global theme is set on the [home page](glossary.html#home-page)

So what happens if you use both `theme` and `global-theme` on
the same page? In the case of the home page, the rule is clear.

Although `global-theme` sets the default theme, it's overriden
by `theme`. This special rule only applies to the home page,
because only the home page lets you set the global theme.

## page theme

A page [theme](glossary.htmle#theme) applies 
theme styling on a per-page basis. It overrides any [global theme](#global-theme)
It's specified using `theme` in the front matter. 

To apply a theme named `informer` to
the currentpage:

```yaml
    ---
    theme: "informer"
    ---
```

### Can a page have both global and page themes?

You may be wondering what happens if you do this:

```yaml
    ---
    global-theme: "pocodocs"
    theme: "informer"
    ---
```

There is only one condition where this makes any kind of sense.

Note this only cano only work on the [home page](glossary-html#home-page), 
because `global-theme` is ignored on all other pages. Since `global-theme`
establishes a preset to be used if no other theme is mentioned on 
a page, and since `theme` overrides the global theme, 
the preceding example would cause the theme named "informer"
to be used on the home page itself. 

PocoCMS has  a [special rule](#special-rule) that the home page normally
uses the global theme, but only if no page theme has been specified.




