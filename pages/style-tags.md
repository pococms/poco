# Style tags

While you can include whole stylesheets, you can also patch changes
with style tags. They're designed to override anything that came
before them, which means they're appended after the stylesheets.
Here's an example that turns all the `<h1>` header text blue:

```yaml
    ---
    style-tags:
    - "article>h1{color:blue}"
    ---
```

