from __future__ import annotations
import web
import internetarchive as ia
from jinja2 import Environment, PackageLoader
import io
import yaml

urls = (
    "/(.*)", "page"
)
app = web.application(urls, globals())
application = app.wsgifunc()

@web.memoize
def get_jinja2_env():
    return Environment(loader=PackageLoader(__name__, "templates"))

def render_template(_filename, **kwargs):
    env = get_jinja2_env()
    template = env.get_template(_filename)
    return template.render(**kwargs)

def get_domain():
    return web.ctx.env.get('HTTP_HOST', '')

def get_itemname():
    domain = get_domain()
    return domain.split(".")[0]

class page:
    def GET(self, path):
        itemname = get_itemname()
        item = self.get_item(itemname)
        return item.render(path)

    def get_item(self, itemname):
        item = Item.find(itemname)
        if not item:
            raise web.notfound("")

        iadata = item.read_archive_yml()
        if iadata and iadata.get("itemtype") == "website":
            root = iadata.get("root") or f"{itemname}.zip"
            return ZipItem(item, root)
        else:
            return item

class ZipItem:
    def __init__(self, item: Item, zip_path: str):
        self.item = item
        self.zip_path = zip_path

    def render(self, path):
        # Internet Archive allows accessing file a/b.txt in file.zip using URL
        # file.zip/a/b.txt
        path = path or "index.html"
        full_path = self.zip_path + "/" + path
        return self.item.read_file(full_path)

CACHE = {}
UNDEFINED = object()

class Item:
    def __init__(self, item: ia.Item):
        self.item = item
        self.files =  [f['name'] for f in item.files]
        self._yamldata = UNDEFINED

    def read_archive_yml(self):
        if self._yamldata is UNDEFINED:
            self._yamldata = self._read_archive_yml()
        return self._yamldata

    def _read_archive_yml(self):
        if "archive.yml" not in self.files:
            return
        contents = self.read_file("archive.yml").decode('utf-8')
        return yaml.safe_load(io.StringIO(contents))

    def get_zipitem(self, path):
        return ZipItem(self, path)

    def has_file(self, path):
        return path in self.files

    def render(self, path):
        if path and not self.has_file(path):
            raise web.notfound("")

        if not path:
            if "index.html" in self.files:
                path = "index.html"
            else:
                return self.render_index()
        return self.read_file(path)

    def render_index(self):
        return render_template(
            "index.html",
            title=self.item.metadata.get('title', ''),
            files=self.item.files)

    def read_file(self, filename) -> bytes:
        file = self.item.get_file(filename)
        response = file.download(return_responses=True)
        return response.content

    @classmethod
    def find(cls, itemname):
        if itemname not in CACHE:
            CACHE[itemname] = cls._find(itemname)
        return CACHE[itemname]

    @classmethod
    def _find(cls, itemname):
        ia_item = ia.get_item(itemname)
        return ia_item and cls(ia_item)
