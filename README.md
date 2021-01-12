# API Stats Interceptor
This is more like watch command in Linux OS but for API endpoints. 

### Compiling from the source files
It’s just like any other Golang project

`go build -o APIStatsInterceptor main.go`

**Note**: You can use a combination of GOOS and GOARCH attributes to compile it for various other operating systems.

### Usage
Assumption:
URL = http://localhost:7700/stats
Headers = `Authorization: <API key>`
Response =
```
{
    "databaseSize": 350871558,
    "lastUpdate": "2021-01-12T02:15:07.501220189Z",
    "indexes": {
        "page_url": {
            "numberOfDocuments": 121909,
            "isIndexing": true,
            "fieldsDistribution": {
                "company_id": 121909,
                "id": 121909,
                "subdomain": 121909,
                "title": 121909
            }
        }
    }
}

```

So to monitor databaseSize, numberOfDocuments, isIndexing, we can form the command as given below,

```
go run main.go -url="http://localhost:7700/stats" -path="DBSize^databaseSize^DATA|#Docs^indexes,page_url,numberOfDocuments^COMMA|Indexing?^indexes,page_url,isIndexing^" -headers=“Authorization:<API key>” -freq=1000
```

Output:
```
Date: Mon, 11 Jan 2021 22:56:18 EST | URL: http://localhost/stats | Frequency: 1s

DBSize: 350 MB
# Docs: 121,909
Indexing?: true
```

Different options:
**url** (required) => can be passed with or without the protocol (http/https)
**path** (required) => values to be monitored with right paths to it along with label and type of units. Each path is separated by _^_ (cap), first value is just a label, second value is the actual path as in JSON response, third value is the unit (COMMA, PERCENT, DATA, P<val> - for prefix, S<val> - for suffix). You can monitor multiple paths and has to be separated by a _|_ (pipe) as in example above.
**freq** (default=1000) => value in milliseconds
**headers** => simple syntax. _key:value_ multiple headers can be passed by structuring it comma separating it. _key:value,key2:value2_