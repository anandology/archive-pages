# Archive Pages

Websites for your archives, hosted directly from your Internet Archive item.

## Overview

How to use:

### Step 1 -- Create your website

```
$ echo "<h1>Hello World!</h1>" > index.html
```

### Step 2 -- Create a zip archive

Zip the entite website and call it archive-pages.zip

```
$ zip -r archive-pages.zip index.html
```

### Step 3 -- Add the zip file to your archive.org item

You can use the [The Internet Archive Python Library][1] to upload.

```
$ ia upload my-item archive-pages.zip
```

[1]:https://archive.org/services/docs/api/internetarchive/cli.html

### Step 4 -- Done!

You website will be live at http://my-item.onarchive.org/
