## apigeecli environments traceconfig update

Update Distributed Trace config for the environment

### Synopsis

Update Distributed Trace config for the environment

```
apigeecli environments traceconfig update [flags]
```

### Options

```
  -p, --endpoint string   Trace endpoint, used only with JAEGER
  -x, --exporter string   Trace exporter can be JAEGER or CLOUD_TRACE
  -h, --help              help for update
  -r, --rate string       Sampler Rate
  -s, --sampler string    Sampler can be set to PROBABILITY or OFF (default "PROBABILITY")
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --default-token    Use Google default application credentials access token
      --disable-check    Disable check for newer versions
  -e, --env string       Apigee environment name
      --metadata-token   Metadata OAuth2 access token
      --no-output        Disable printing all statements to stdout
  -o, --org string       Apigee organization name
      --print-output     Control printing of info log statements (default true)
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli environments traceconfig](apigeecli_environments_traceconfig.md)	 - Manage Distributed Trace config for the environment

###### Auto generated by spf13/cobra on 19-Dec-2023
