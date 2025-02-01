# mcphost
Model Context Protocol server to share files/directories, prompts and resources from the host machine.


# JSON RPC requests

```
{"method":"tools/list","params":{},"jsonrpc":"2.0","id":1}

{"method":"tools/call","params":{"name":"notion_get_page","arguments":{"pageId":"1406f674ce1e80bc90c0c9d1814b4bb0"}},"jsonrpc":"2.0","id":16}

{"method":"tools/call","params":{"name":"ping","arguments":{"message":"hello"}},"jsonrpc":"2.0","id":28}


{"method":"prompts/list","params":{},"jsonrpc":"2.0","id":17}

{ "jsonrpc": "2.0", "id": 2, "method": "prompts/get", "params": { "name": "hello", "arguments": { "name": "Pierre Carion" } } }

{"method":"ping","params":{},"jsonrpc":"2.0","id":42}


```