## apigeecli datastores list

Returns a list of data stores in the org

### Synopsis

Returns a list of data stores in the org

```
apigeecli datastores list [flags]
```

### Options

```
  -h, --help                help for list
      --targettype string   TargetType is used to fetch datastores of matching type; default is fetch all
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

* [apigeecli datastores](apigeecli_datastores.md)	 - Manage data stores to export analytics data

###### Auto generated by spf13/cobra on 19-Dec-2023
