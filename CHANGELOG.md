# Change Log

## Unreleased

## 0.3.0 -- 2020-09-06

* The sytem now expects a single zip file with name archive-paghes.zip to be present in the item
* The archive.yml is depricated

## 0.2.0 -- 2020-08-24

* Switched the implementation from Python to Go lang
* Added docker-compose to streamline development
* Added redis for caching metadata and Varnish for http caching
* Only the items that have archive.yml will be supported

## 0.1.0 -- 2020-07-27

* First working prototype of archive pages
