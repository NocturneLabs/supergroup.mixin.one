# Mixin Supergroup Plugins

## Getting Started

Due to the way go plugins works, all plugins must build with same source of host program. 
In other word, the directory layout must be:

```
.
├── supergroup.mixin.one
│   ├── config
│   ├── durable
│   ├── gears
│   ├── interceptors
│   ├── middlewares
│   ├── models
│   ├── plugin
│   ├── plugins
│   ├── routes
│   ├── services
│   ├── session
│   ├── views
│   └── web
└── supergroup.mixin.one-plugins
    ├── build
    └── course_record
```

Let's build our plugins!

```bash
make
```

After running commands above, all plugins will be built to shared libraries located at `build/`.
Then update host program's config to load our plugins.

```yaml
plugins:
  - shared_library: ../supergroup.mixin.one-plugins/build/course_record.so
    config:
      database_url: postgres://qinix:@localhost/course_record_dev?sslmode=disable
```

...And, run host program as usual.

```bash
# current working directory: supergroup.mixin.one/
go run .
```

## Plugins

### course_record

The course_record plugin will record all messages when group prohibited status is on.
Admin can publish records manually.

#### Config

```yaml
    config:
      database_url: postgres://username:password@host:port/dbname?sslmode=disable
```

#### API

Note that all endpoints below require authenticated.

##### GET /course_record/records

List all records. Return public records only if user is not admin.

Response:
```json
[
  {
    "id": 1,
    "title": "Untitled",
    "is_public": false,
    "started_at": "2019-07-11T17:17:21.10851+08:00",
    "finished_at": "2019-07-11T17:18:34.406684+08:00",
    "speakers": [
      {
        "id": "5d8e3dda-e80b-4f6b-bd74-e5cfd3947ddd",
        "full_name": "章琦",
        "avatar_url": "https://images.mixin.one/4il_Op2IgGcB4HGGgi5S0Cri7-9LAz7ml6UR27zlPkQwFE9HYru2sUbj6UFegBGtl7mRT9cRT5dPJkxS8RT3AEw=s256"
      }
    ]
  },
  {
    "id": 2,
    "title": "Untitled",
    "is_public": false,
    "started_at": "2019-07-11T17:18:34.408431+08:00",
    "finished_at": "2019-07-11T17:19:03.313064+08:00",
    "speakers": [
      {
        "id": "5d8e3dda-e80b-4f6b-bd74-e5cfd3947ddd",
        "full_name": "章琦",
        "avatar_url": "https://images.mixin.one/4il_Op2IgGcB4HGGgi5S0Cri7-9LAz7ml6UR27zlPkQwFE9HYru2sUbj6UFegBGtl7mRT9cRT5dPJkxS8RT3AEw=s256"
      }
    ]
  }
]
```

##### PUT/PATCH /course_record/records/:record_id

Update record's title or is_public

Params:
- title: string (optional)
- is_public: string (optional)

Response:
```json
{
  "id": 2,
  "title": "test",
  "is_public": true,
  "started_at": "2019-07-11T17:18:34.408431+08:00",
  "finished_at": "2019-07-11T17:19:03.313064+08:00",
  "speakers": [
    {
      "id": "5d8e3dda-e80b-4f6b-bd74-e5cfd3947ddd",
      "full_name": "章琦",
      "avatar_url": "https://images.mixin.one/4il_Op2IgGcB4HGGgi5S0Cri7-9LAz7ml6UR27zlPkQwFE9HYru2sUbj6UFegBGtl7mRT9cRT5dPJkxS8RT3AEw=s256"
    }
  ]
}
```

##### GET /course_record/records/:record_id/messages

List record's messages

Response:
```json
[
  {
    "message_id": "73ab069a-a070-4d7c-a3e7-f114a52fce0b",
    "quote_message_id": "",
    "data": "dGVzdA==",
    "category": "PLAIN_TEXT",
    "created_at": "2019-07-11T17:17:35.133223+08:00",
    "speaker": {
      "id": "5d8e3dda-e80b-4f6b-bd74-e5cfd3947ddd",
      "full_name": "章琦",
      "avatar_url": "https://images.mixin.one/4il_Op2IgGcB4HGGgi5S0Cri7-9LAz7ml6UR27zlPkQwFE9HYru2sUbj6UFegBGtl7mRT9cRT5dPJkxS8RT3AEw=s256"
    }
  },
  {
    "message_id": "6ed38fb3-1f23-4a7f-8a37-8943ae749f40",
    "quote_message_id": "",
    "data": "YWFh",
    "category": "PLAIN_TEXT",
    "created_at": "2019-07-11T17:17:36.135853+08:00",
    "speaker": {
      "id": "5d8e3dda-e80b-4f6b-bd74-e5cfd3947ddd",
      "full_name": "章琦",
      "avatar_url": "https://images.mixin.one/4il_Op2IgGcB4HGGgi5S0Cri7-9LAz7ml6UR27zlPkQwFE9HYru2sUbj6UFegBGtl7mRT9cRT5dPJkxS8RT3AEw=s256"
    }
  }
]
```
