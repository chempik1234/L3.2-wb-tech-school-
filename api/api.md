# TODO: api.md

## До вступления

> **WHY NOT openAPI?**

How do I even describe a redirect in it?

## Paths

1. **POST /shorten** - Create link

* Input:

```json
{
  "source_url": "https://ya.ru"
}

// for custom short link
{
  "source_url": "https://ya.ru",
  "short_url": "1234"
}

// Warning: custom link ignored if empty
{
  "source_url": "https://ya.ru",
  "short_url": ""
}
```

* Output:

```json
{
  "source_url": "https://ya.ru",
  "short_url": "ksola",
  "created_at": 1000000000
}
```

* Validation: **short_url** must either be null or have <=30 chars & be **unique**

---

2. **GET /s/{short_url}** - Redirect to Short URL

* Output: redirect to the original URL.
* Validation: **short_url** must exist; otherwise 404.

---

3. **GET /analytics/{short_url}** - Analytics for Short URL

* Input: None
* Output:
```json
  {
  "source_url": "https://ya.ru",
  "short_url": "<short_url>",
  "total_redirects": 2,
  "data": [
    {
      "click_at": 1000000000,
      "user_agent": "..."
    },
    {
      "click_at": 1000000001,
      "user_agent": "..."
    }
  ]
}
```

* Validation: **short_url** must exist; otherwise 404.