# PocoCMS release checklist

* go fmt
* Run all stylesheets through `vnu`
* `tree` on the source directory
* `du -h` on the source directory
* Check document Markdown source code for ocurrences of TODO:, 
including as comments in Go template format
* Run poco in the poco directory
* Run this script. ASSUMES ~/pococms/poco/foobar/ IS A DISPOSABLE DIRECTORY
```bash
#!/bin/zsh
TESTDIR=~/pococms/poco/foobar/
cd ~/pococms/poco/
go build
rm -rf $TESTDIR
mkdir -p $TESTDIR
cd $TESTDIR
poco
ls $TESTDIR/.poco/themes/informer
```
## Interactive theme tests
* For each theme click all links in all page layout elements

