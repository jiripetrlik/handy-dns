### CURL examples

```
# List items
curl http://localhost:8080/api/list

# Add item
curl "http://localhost:8080/api/create?name=machine1&class=IN&itemType=A&data=127.0.0.1"

# Update item
curl "http://localhost:8080/api/update?id=1&name=machine5&class=IN&itemType=A&data=127.0.0.1"

# Delete item
curl "http://localhost:8080/api/delete?id=1"
```
