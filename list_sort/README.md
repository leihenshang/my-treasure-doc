# list sort

## description

used to list sort fields.

web:
```json
{
    "page": 1,
    "pageSize": 10,
    "sort": "[{ \"field\": \"id\", \"order\": \"asc\" },{ \"field\": \"name\", \"order\": \"desc\" }]"
    // or 
    "sort": "id_asc,name_desc"
}

```
## plan

### step 1
 - [x] support json sort format.
 - [x] support string sort format.
 ### step 2
 - [ ] support sort field and order with default value.
 - [ ] add options mode