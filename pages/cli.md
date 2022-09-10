---
theme: ".poco/themes/probot"
---
# PocoCMS command-line options

{{- /* TODO: screenshots of both kinds of output for timestamp */ -}}

## timestamp

The `-timestamp` command-line option lets you insert
a timestamp before the rest of the article text on your
home page. Useful when server or browser caches are
preventing you from changes to your site's output.

### timestamp example

Enter this at the operating system prompt:

```bash
poco -timestamp
```

Upon completion you'll see the timestamp before the completion message:

```
25 Aug 2023 18:33:29 Site published to WWW/index.html
```

You can now check it against the timestamp on your homepage
to ensure the home page you see has cleared browser or server caches.


## webroot

The `-webroot` command-line option lets you choose what directory
Poco CMS uses for the finished HTML, CSS, images, and related assets
to be published. By default it's `WWW`. 

To change it, just run `poco` `-webroot` followed by the
directory you'd like to use. For example:

```
$ poco -webroot  public_html

```


You'll see results something like this:

```
PocoCMS: Site published to /Users/travis/mysite/public_html
```
