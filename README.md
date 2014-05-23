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
page directory. The contents of the draft file are ignored.
