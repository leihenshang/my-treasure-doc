

### CommonJS

Then configure in your entry file:

```javascript
import Vue from 'Vue'
import VueAxiosPlugin from 'vue-axios-plugin'

Vue.use(VueAxiosPlugin, {
  // request interceptor handler
  reqHandleFunc: config => config,
  reqErrorFunc: error => Promise.reject(error),
  // response interceptor handler
  resHandleFunc: response => response,
  resErrorFunc: error => Promise.reject(error)
})
```

## Options

Except axios default [request options](https://github.com/axios/axios#request-config), `vue-axios-plugin` provide below request/response interceptors options:

|Name|Type|Default|Description|
|:--:|:--:|:-----:|:----------|
|**[`reqHandleFunc`](#)**|`{Function}`|`config => config`|The handler function for request, before request is sent|
|**[`reqErrorFunc`](#)**|`{Function}`|`error => Promise.reject(error)`|The error function for request error|
|**[`resHandleFunc`](#)**|`{Function}`|`response => response`|The handler function for response data|
|**[`resErrorFunc`](#)**|`{Function}`|`error => Promise.reject(error)`| The error function for response error |

## Example

Default method in `$http`, it just contains get and post method:

```javascript
this.$http.get(url, data, options).then((response) => {
  console.log(response)
})
this.$http.post(url, data, options).then((response) => {
  console.log(response)
})
```

Use axios original method in `$axios`, by this, you can use all allowed http methods: get,post,delete,put...

```javascript
this.$axios.get(url, data, options).then((response) => {
  console.log(response)
})

this.$axios.post(url, data, options).then((response) => {
  console.log(response)
})
```

## TODO

- [ ] Unit test.

## Notice!!!

When you send a request use `application/x-www-form-urlencoded` format, you need to use [qs](https://github.com/ljharb/qs) library to transform post data, like below:

```js
import qs from 'qs'
this.$http.post(url, qs.stringify(data), {
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
  }
}).then((response) => {
  console.log(response)
})
```

But if the `data` has properties who's type if `object/array`, you need convert these properties into JSON string:

```js
import qs from 'qs'

function jsonProp (obj) {
  // type check
  if (!obj || (typeof obj !== 'object')) {
    return obj
  }
  Object.keys(obj).forEach((key) => {
    if ((typeof obj[key]) === 'object') {
      obj[key] = JSON.stringify(obj[key])
    }
  })
  return obj
}

this.$http.post(url, qs.stringify(data), {
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
  },
  transformRequest: [
    function (data) {
      // if data has object type properties, need JSON.stringify them.
      return qs.stringify(jsonProp(data))
    }
  ]
}).then((response) => {
  console.log(response)
})
```
