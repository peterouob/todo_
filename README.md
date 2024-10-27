## api一覽

|           | query | token | method | database    | group |
|-----------|-------|-------|--------|-------------|-------|
| /login    | no    | no    | post   | mysql/redis | no    |
| /register | no    | no    | post   | mysql       | no    |
| /         | no    | yes   | get    | mongo/redis | todo  |
| /:id      | param | yes   | get    | mongo/redis | todo  |
| /:id      | param | yes   | put    | mongo/redis | todo  |
| /:id      | param | yes   | delete | mongo/redis | todo  |
| /filter   | query | yes   | get    | mongo/redis | todo  |
| /done     | no    | yes   | post   | mongo/redis | todo  |
| /done/:id | param | yes   | get    | mongo/redis | todo  |
| /create   | no    | yes   | post   | mongo/redis | todo  |
