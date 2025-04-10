## apigeecli products create

Create an API product

### Synopsis

Create an API product

```
apigeecli products create [flags]
```

### Options

```
  -f, --approval string        Approval type
      --attrs stringToString   Custom attributes (default [])
  -d, --desc string            Description for the API Product
  -m, --displayname string     Display Name of the API Product
  -e, --envs stringArray       Environments to enable
      --gqlopgrp string        File containing GraphQL Operation Group JSON. See samples for how to create the file
      --grpcopgrp string       File containing gRPC Operation Group JSON. See samples for how to create the file
  -h, --help                   help for create
  -i, --interval string        Quota Interval
  -n, --name string            Name of the API Product
      --opgrp string           File containing Operation Group JSON. See samples for how to create the file
  -p, --proxies stringArray    API Proxies in product. Ex: -p proxy1 -p proxy2
  -q, --quota string           Quota Amount
  -s, --scopes stringArray     OAuth scopes
  -u, --unit string            Quota Unit
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

* [apigeecli products](apigeecli_products.md)	 - Manage Apigee API products

###### Auto generated by spf13/cobra on 19-Dec-2023
