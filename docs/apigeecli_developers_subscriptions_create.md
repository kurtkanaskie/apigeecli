## apigeecli developers subscriptions create

Create a subcription for a developer

### Synopsis

Create a developer

```
apigeecli developers subscriptions create [flags]
```

### Options

```
  -p, --apiproduct string   The name of the API Product
  -e, --email string        The developer's email
  -d, --end string          Subscription end time
  -h, --help                help for create
  -n, --name string         The subscription name
  -s, --start string        Subscription start time
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

* [apigeecli developers subscriptions](apigeecli_developers_subscriptions.md)	 - Manage subscriptions for a Developer

###### Auto generated by spf13/cobra on 19-Dec-2023
