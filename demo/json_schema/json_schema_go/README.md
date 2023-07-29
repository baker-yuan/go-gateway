

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "age": {
      "type": "integer",
      "minimum": 0
    },
    "email": {
      "type": "string",
      "format": "email"
    }
  },
  "required": ["name", "age"]
}
```





```json
{
  type: 'object',
  required: [
    'userName',
    'age',
  ],
  properties: {
    userName: {
      type: 'string',
      title: '用户名',
      default: 'Liu.Jun',
    },
    age: {
      type: 'number',
      title: '年龄'
    },
    bio: {
      type: 'string',
      title: '签名',
      minLength: 10,
      default: '知道的越多、就知道的越少',
      'ui:options': {
        placeholder: '请输入你的签名',
        type: 'textarea',
        rows: 1
      }
    }
  }
}
```

