# poco command-line options

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
