# Building PocoCMS from source

Until 15 September 2022 PocoCMS must be built from source using Go.
It's easy to install Go and even to build, because
PocoCMS is a single file.

### One time: install Go and git

* Install the [Go language](https://go.dev/dl/) if necessary.
* Install [git](https://git-scm.com/downloads) if necessary.
* Create a directory or change to a directory to install PocoCMS.

```bash
mkdir ~/pococms
```

* Navigate to that directory.

```bash
cd ~/pococms
```

* Clone the PocoCMS repo

```bash
git clone https://github.com/pococms/poco
```

* The repo is now in ~/pococms/poco, (in this example) so navigate there.

```bash
cd poco
```

### One time: compile PocoCMS

* And compile: 

```bash
go build 
```

**OR...**


* There's only one file, so you can also use go run.
That runs the go compiler on the single `main.go` 
containing PocoCMS, then executes PocoCMS.

```bash
go run main.go
```


Add the `poco` executable to your system path 
so you can run it from any directory.

### PocoCMS creates a starting project automatically.


