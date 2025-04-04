## apigeecli apis fetch

Returns a zip-formatted proxy bundle 

### Synopsis

Returns a zip-formatted proxy bundle of code and config files

```
apigeecli apis fetch [flags]
```

### Options

```
  -h, --help          help for fetch
  -n, --name string   API Proxy Bundle Name
  -v, --rev int       API Proxy revision. If not set, the highest revision is used (default -1)
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --default-token    Use Google default application credentials access token
      --disable-check    Disable check for newer versions
      --metadata-token   Metadata OAuth2 access token
      --no-output        Disable printing all statements to stdout
  -o, --org string       Apigee organization name
      --print-output     Control printing of info log statements (default true)
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli apis](apigeecli_apis.md)	 - Manage Apigee API proxies in an org

###### Auto generated by spf13/cobra on 19-Dec-2023
