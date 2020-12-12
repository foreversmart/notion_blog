# notion_blog
this is a tool which transfer notion page to hugo html page

### Install
go get github.com/foreversmart/notion_blog

### Usage
./export it will generate hugo html page to target path usually hugo_blog/posts/

### Config
``` json
{
  "page_ids": [
    "aa6f1f61bd114644a4c309cc31895d5e",
    "Istio-f76464d369804e42bb67b180e7155a11",
    "https://www.notion.so/Golang-79eabfceb5cb4256a341c7cc23a2264a",
    "https://www.notion.so/go-panic-and-recover-676eb9f86a5d4e43b35647bec3dd04ac",
    "https://www.notion.so/Nginx-request-body-a0d919fd7550427083ca36ed9fcc7ca6",
    "https://www.notion.so/20ecc4449a3843e6bea9324c3328241e",
    "https://www.notion.so/635122dd9eb24c33ac08b3e71f1a7fbd"
  ],
  "out_put_path": "/Users/someone/hugo_blog/content/post"
}
``` 
* page_ids is an array of notion public page id which you want generate these to hugo html page
* out_put_path define the target path export tool will generate it into

### Notion Page Meta
you can redefine some notion page meta by yourself instead use tool generated.
notion blog tool craw notion page comment as notion page meta which start with **meta:**
, the same as tags by **tags:**, and categories as **cate:**

* Meta Properties List
```
tags:$tags
cate:$categories
meta:title:$title
meta:created_at:$created_at
meta:sub_title:sub_title
meta:desc:desc
meta:url_path:url_path
meta:url:url
```
