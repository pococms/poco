# PocoCMS Frequently Asked Questions

How do I get blank lines after paragraphs? 
: To get end a line and start a new one, you need
to end the first line with two spaces. Imagine the
little dot characters ·· are actually you pressing
the space key twice:


```
Line 1.··   ← Pretend those are 2 space characters
Line 2.
```

**What? Why? This is terrible!**

It's actually the way HTML behaves, and remember that
Markdown was made to be easy to read by itself, but
to [translate directly to HTML](https://daringfireball.net/projects/markdown/) when you did a conversion.

*Some background on HTML and how it works with Markdown*

And in HTML, any number of spaces and/or paragraphs (called `newlines` in the programming world) between words are automatically
replaced by a single space. So, if you write any of these
in Markdown:


```
Line 1.
Line 2.
```

Or this:

```
Line 1.


Line 2.
```

Or even this:

```
Line 1.



Line 2.
```

The result will always be this:

```
Line 1. Line 2.
```

In HTML, if you write this:

```
<p>Line 1.<br/ >Line 1.</p>
``` 

you will get the result you're looking for:

```
Line 1. 
Line 2.
```

There may or may not be a space between paragraphs
depending on how the style sheet was written.

You will get that same result if your code looks like this:

```
<p>Line 1.</p>



<p>Line 1.</p>
``` 

It will still appear this way on the web page:

```
Line 1. 

Line 2.
```

Adding 2 spaces to the end of a line is how to get Markdown
to enclose each line in paragraph tags, but to use a line break between them. Again, if
you type 2 spaces at the end of the first line, 
both will be rendered in HTML as paragraphs. 
This Markdown:

```
Line 1.··   ← Pretend those are 2 space characters
Line 2.
```

Becomes this HTML:

```html
<p>Line 1.<br/ >Line 2.</p>
``` 

And will render something like this:


```
Line 1. 
Line 2.
```
