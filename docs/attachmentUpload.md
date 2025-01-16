1. User enters Message
2. Clicks on send
3. client -> server `HTTP POST /messages`
  ``` json
  {
    "content": "Hello",
    "attachments": [
      {
        "filename": "image.jpg",
        "type": "image/jpeg",
      },
      {
        "filename": "main.py",
        "type": "text/plain",
      }
    ]
  }
  ```
4. server -> client `201 Created`
  ``` json
  {
    "id": "123",
    "author": "2345",
    "content": "Hello",
    "attachments": [
      {
        "id": "345",
        "filename": "image.jpg",
        "type": "image/jpeg",
        "url": "https://example.com/image.jpg",
      },
      {
        "id": "456",
        "filename": "main.py",
        "type": "text/plain",
        "url": "https://example.com/main.py",
      }
    ]
  }
  ```
5. client -> server `HTTP POST /attachments/345`
  ``` blob
  sadfkjshfdkljhdfsg
  ```
1. server -> client `201 Created`
2. client -> server `HTTP POST /attachments/456`
  ``` blob
  sadfkjshfdkljhdfsg
  ```
1. server -> client `201 Created`