---
pagetheme: mytheme 
hide: nav
aside: aside.md
sidebar: left
---
##### Blog

# Static Site Generators


#### By **Anna Utopia** | *December 3*   [![LinkedIn](../.poco/img/linkedin-24px-blue-outline.svg)](https://linkedin.com/) [![YouTube logo](../.poco/img/youtube-24px-red-outline.svg)](https://youtube.com/@pococms/)   [![Facebook logo](../.poco/img/facebook-24px-blue-outline.svg)](https://facebook.com/)   [![Instagram logo](../.poco/img/instagram-24px-magenta-outline.svg)](https://www.instagram.com/e.emerald.repair/)  [![Twitter logo](../.poco/img/twitter-24px-blue-outline.svg)](https://www.instagram.com/e.emerald.repair/)


As you may have noted from my bio, it's almost like I was born to
promote the idea of static site generators. If the universe actually
were a simulation, I would be a character designed solely to move
the idea of static site generators forward.


## Spoiler alert: PocoCMS

There are plenty of good candidates out there. Let's start with my favorite, which many people know is [PocoCMS](https://pococms.com).
Off the top of my head, I like it because:

* PocoCMS is free and open source. You can use it to make money and not pay
the PocoCMS people anything.
* It has a good set of site templates to start out with.
* The sites it generates take very little time to load. People are
way less likely to skip to another site.
* You can build whole templates using Markdown!
* You can control whether the header, navbar, sidebar, and footer are
displayed  on any page, at least if it's created with their theme framework.

### Fastest way to get on the web?

I think PocoCMS is the fastest of the static site generators to get
a full, professional-looking site on the web. Watch:

* Download PocoCMS
* Create your site:

```
poco -new annautopia
cd annautopia
```

There's already a site! But let's change its theme.

```
# Replace edit with whatever editor you use
edit index.md
```

When you get `index.md`, set a global theme for the site:

##### Filename: **index.md**

```
---
theme: pocodocs
---

# Welcome to annautopia

Hit up my [Insta](instagram.com/)!
```









