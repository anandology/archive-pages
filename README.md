# onarchive

Onarchive is a tool to serve an archive.org item into a website.

## How it works

If you have an archive.org item `foobar`, you can see it at `foobar.onarchive.org`.

By default, it renders the website based on the type of the item, but it can be customized by writing an `archive.yml` file.

## Hosting a website

One of the most common use cases of onarchive is to host an archive of a website. To do that, create a zipfile of the website and add it to the item.

For example, if name of the item is mygovwebsite2020, add a zipfile mygovwebsite2020.zip with all the contents of the archive.

Add an archive.yml file to the item with the following contents:

```
version: 1
itemtype: website
root: mygovwebsite2020.zip
```

Once that is done, the website will be available at <https://mygovwebsite2020.onarchive.org>

