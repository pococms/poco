# Front Matter

# TODO: SkipPublish isn't implemented yet

## skippublish

`SkipPublish` lists files and directories you don't want to be published.

Remember that if a directory contains the files `index.md`, `installation.md`,
and `avatar.png`, and `401K-info.xls`, here's what will happen when you 
run PocoCMS:

* `index.md` and `installation.md` will be converted to HTML document files
named `index.html` and `installation.html`. The the HTML files will be published
The Markdown files will not.
* By designed, all files other than than Markdown files get published. That's
because they're assumed to be required for the site. In this example,
it makes sense that `avatar.png` gets published but you may not want
the personal spreadsheet with your 401K details published.
* Likewise, you probably don't want directories like `node_modules`,
`.DS_Store`, or `.git` published. 

`.git` is included here for good form but 
PocoCMS treats directories with that name starts with `.` as hidden and
doesn't publish them.
* The answer to these problems is to list what you don't want published in `SkipPublish` 
as shown below.

### Example
      ---
      SkipPublish:
        - 401k-info.xls
        - node_modules
        - www
        - .git
        - .DS_Store
        - .gitignore
      ---

Here are some common items to skip:

      ---
      SkipPublish:
        - node_modules
        - htdocs
        - public_html
        - www
        - .git
        - .DS_Store
        - .gitignore
      ---



