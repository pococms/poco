# PocoCMS release checklist

* Change to the poco directory
* go build
* go test
* Check quickly through the source code
* go fmt
* Run all stylesheets through `vnu` via my `v` utility
* `du -h` on the source directory
* `tree` on the source directory
* Check document Markdown source code for ocurrences of TODO:, 
including as comments in Go template format
* Run poco in the poco directory
* Run the test script `/.ts1` in the poco directory
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

