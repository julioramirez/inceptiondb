# Retrieve collection

Curl example:

```sh
curl "https://example.com/v1/collections/my-collection" \
-d '{
    "name": "my-collection"
}'
```


HTTP request/response example:

```http
GET /v1/collections/my-collection HTTP/1.1
Host: example.com

{
    "name": "my-collection"
}

HTTP/1.1 200 OK
Content-Length: 47
Content-Type: text/plain; charset=utf-8
Date: Mon, 15 Aug 2022 02:08:13 GMT

{
    "indexes": 0,
    "name": "my-collection",
    "total": 0
}
```


