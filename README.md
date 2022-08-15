---
Header: ./pages/diagnostics/layout/header.md
Nav: ./pages/diagnostics/layout/nav.md
Aside: ./pages/diagnostics/layout/aside.md
Footer: ./pages/diagnostics/layout/footer.md

StyleFiles:
  - pages/assets/css/poquito.css
  #- pages/assets/css/new-poquito.css
---
## Diagnostics

Style sheets: {{ .StyleFiles }}

## Sidebar tests

### NEWER left and right sidebars

* **Left sidebar using newer-poquito.css** [nwr-poquito-lside](nwr-poquito-lside.html)
* **Right sidebar using newer-poquito.css** [nwr-poquito-lside](nwr-poquito-rside.html)

### NEW left and right sidebars

[new-left-sidebar](new-left-sidebar.html)

[new-right-sidebar](new-right-sidebar.html)

[Probot right sidebar](probot.html) 

[Probot left sidebar](probot-left.html) theme

## OLD left and right sidebars
`poquito.css` has left and right sidebar support. Right sidebars are
supported by default. To get sidebar support:

* Add a Markdown file for `Aside` in the
front matter as shown below (it can be any Markdown file and doesn't
have to be named `aside`).
* Add `poquito.css` under `StyleFiles`:

      ---
      Aside: aside.md
      StyleFiles:
        - poquito.css
      ---

Here are minimal examples of left and right sidebars:

[Sidebar right example](sidebar-right.html)

[Sidebar left example](sidebar-left.html)

## I can't see ANY sidebars!

If you can't see the sidebars you're using a mobile device. The `poquito.css`
media queries hide the sidebar, header, nav, and footer.

## Long text

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Sagittis purus sit amet volutpat consequat mauris nunc congue nisi. Porttitor massa id neque aliquam vestibulum morbi. Non blandit massa enim nec dui nunc mattis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras. At consectetur lorem donec massa. Id cursus metus aliquam eleifend mi in nulla. Blandit cursus risus at ultrices. Vel risus commodo viverra maecenas accumsan. Urna porttitor rhoncus dolor purus non enim praesent elementum. Nec tincidunt praesent semper feugiat nibh sed pulvinar proin. Dignissim sodales ut eu sem integer. Facilisis mauris sit amet massa vitae tortor condimentum lacinia quis.


Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Sagittis purus sit amet volutpat consequat mauris nunc congue nisi. Porttitor massa id neque aliquam vestibulum morbi. Non blandit massa enim nec dui nunc mattis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras. At consectetur lorem donec massa. Id cursus metus aliquam eleifend mi in nulla. Blandit cursus risus at ultrices. Vel risus commodo viverra maecenas accumsan. Urna porttitor rhoncus dolor purus non enim praesent elementum. Nec tincidunt praesent semper feugiat nibh sed pulvinar proin. Dignissim sodales ut eu sem integer. Facilisis mauris sit amet massa vitae tortor condimentum lacinia quis.

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Sagittis purus sit amet volutpat consequat mauris nunc congue nisi. Porttitor massa id neque aliquam vestibulum morbi. Non blandit massa enim nec dui nunc mattis. Vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras. At consectetur lorem donec massa. Id cursus metus aliquam eleifend mi in nulla. Blandit cursus risus at ultrices. Vel risus commodo viverra maecenas accumsan. Urna porttitor rhoncus dolor purus non enim praesent elementum. Nec tincidunt praesent semper feugiat nibh sed pulvinar proin. Dignissim sodales ut eu sem integer. Facilisis mauris sit amet massa vitae tortor condimentum lacinia quis.



