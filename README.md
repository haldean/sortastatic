sortastatic
===========

This is another attempt at creating a static site server that I like. This works
on a sorta-blog-like structure that more page-oriented than time-oriented. Pages
are created in an "index" directory that contains a subdirectory for each page.
The name of each subdirectory is the URL that that page will be accessible at.
Inside each subdirectory should be one markdown file and any number of other
files; when a user goes to `yoursite.com/x/yourpage/`, the markdown file at
`index-directory/yourpage/file.md` will be loaded, converted to HTML and served.
All files within the `/x/yourpage/` URL path will be served as static files from
the page directory.

There is also support for site-wide common files, stored in the "common"
directory and accessed as `yoursite.com/c/whatever.png`. This is where
Javascript, CSS, etc. that are used on every page in the site can be placed.

Go templates are used to render the pages; the structure passed to the page
template is:

      type Page struct {
        Name     string
        Path     string
        Title    string
        Body     string
        Url      string
        Rendered string
        Public   bool
      }

The structure passed to the index page is a list of page structs sorted by
title.

If you want to prevent a page from showing up on the index but still want it to
be accessible through direct links, you can place a file named "draft" in its
page directory. The contents of the draft file are ignored. To ignore a
directory entirely (no page at all, not even by URL), place an empty file called
"ignore" in the directory.

sortastatic makes heavy use of caching; pages are enumerated and templates are
loaded at start time, and Markdown content is lazily rendered, then cached.
If you want to force a cache flush without restarting the process (which will
reload templates and the page list, and clear all cached rendered content), send
the sortastatic process a SIGUSR1 signal (this only works on Unices, sorry
Windows users). There's a convenience utility built into sortastatic to do this;
just run `sortastatic -reload`.

## Requirements

* [github.com/russross/blackfriday](github.com/russross/blackfriday)
  Really easy-to-use Markdown library. &lt;3.
