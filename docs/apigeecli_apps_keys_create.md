## apigeecli apps keys create

Create a developer app key

### Synopsis

Create a a developer app key

```
apigeecli apps keys create [flags]
```

### Options

```
      --attrs stringToString   Custom attributes (default [])
  -x, --expires string         A setting, in seconds, for the lifetime of the consumer key
  -h, --help                   help for create
  -k, --key string             Developer app consumer key
  -p, --prods stringArray      A list of api products
  -s, --scopes stringArray     OAuth scopes
  -r, --secret string          Developer app consumer secret
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --default-token    Use Google default application credentials access token
  -d, --dev string       Developer email
      --disable-check    Disable check for newer versions
      --metadata-token   Metadata OAuth2 access token
  -n, --name string      Developer App name
      --no-output        Disable printing all statements to stdout
  -o, --org string       Apigee organization name
      --print-output     Control printing of info log statements (default true)
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli apps keys](apigeecli_apps_keys.md)	 - Manage developer app keys

###### Auto generated by spf13/cobra on 19-Dec-2023
