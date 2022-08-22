---
Title: "pocoCMS: Display all known information about this document"
Description: "PocoCMS: Markdown-based CMS in 1 file, written in Go"
Theme: SUPPRESS
Author: "Tom Campbell"
Keywords: "static site generator, CMS, wordpress replacement, Markdown"
LinkTags:
    - <link rel="icon" href="favicon.ico">
    - <link rel="preconnect" href="https://fonts.googleapis.com">
    - <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    - <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&display=swap" rel="stylesheet">
#StyleFiles: 
#    - "https://cdn.jsdelivr.net/gh/pococms/poco/pages/assets/css/poquito.css"
#    - 'https://cdn.jsdelivr.net/npm/holiday.css'
SkipPublish:
    - node_modules
    - htdocs
    - public_html
    - WWW
    - .git
    - .DS_Store
    - .gitignore
---
# {{ .Title }}

| Setting               | Value                      |
| --------------------- | -------------------------- |
| Title                 | {{ .Title }}               |
| Description           | {{ .Description }}         |
| Author                | {{ .Author }}              |
| Keywords              | {{ .Keywords }}            |
| Stylesheets           | {{ .StyleFiles }}          |
