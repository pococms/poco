#!/bin/zsh
# v Uses the vnu utility to check the syntax of a .CSS file.
# vnu only knows about HTML files. This takes the CSS file,
# inserts its contents into a dummy HTML file between
# <style> tags, places that file into a temporary directory, 
# and runs vnu on the temp HTML file instead of the original CSS file.
#
# To obtain vnu, see:
# https://github.com/validator/validator/releases

# Ensure a file was named on the command line.
[ ! -f $1 ] && echo "Please specify a file" &&  exit 1

# Ensure vnu is available.
if ! command -v vnu &> /dev/null
then
    echo "vnu needs to be on the path. Try visiting:"
    echo "https://github.com/validator/validator/releases"
    exit 1
fi

# Create a temporary directory.
# Supposed to work in both MacOS and Linux. Thank you, timo-tijhof!
# https://unix.stackexchange.com/questions/30091/fix-or-alternative-for-mktemp-in-os-x
tmpdir=$(mktemp -d 2>/dev/null || mktemp -d -t 'tmpdir')

# Read the contents of the .CSS file named on the command line as $1
# into the variable $infile
infile=$(<$1)

# Insert the .CSS file into a one-line HTML file.
# That way if there are errors, the correct line numbers
# will be reported.
read -r -d '' contents << EOM
<!DOCTYPE html><html lang="en"><head><title>$1</title><style>$infile</style></head><body></body></html>
EOM

# Extract the name only of the CSS file (no directory)
basename="${1##*/}"

# Append that filename to the temporary directory.
outfile=$tmpdir/$basename 

# Create a file in the temporary directory. It's the one-line
# HTML file with the .CSS file embedded.
echo "$contents" > $outfile

# Check its syntax.
# at some point extract the filename in quotes
# using a regexp something like ^"(.*?)"
vnu $outfile  > foo.1

if [ $? -ne 1 ]; then
  echo "No errors found in $1"
fi




