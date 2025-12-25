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
  "created_at": "...iso datetime"
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
  "unique_user_agents": 2,
  "data": [
    {
      "minute": "...iso datetime",
      "clicks_in_minute": 120469101,
      "data": [
        {
          "user_agent": "...",
          "clicks": 1
        }
      ]
    },
    {
      "minute": "...prev datetime + 2 minutes",
      "clicks_in_minute": 1498532,
      "data": [
        {
          "user_agent": "...",
          "clicks": 1
        }
      ]
    }
  ]
}
```

* Validation: **short_url** must exist; otherwise 404.