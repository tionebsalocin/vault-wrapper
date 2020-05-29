# vault-wrapper
This wrapper uses kerberos to authenticate to Vault, collect secrets and pass those to a command it will execute through environment variables.

The list of secret to collect is obtained by setting evironment variables with following syntax:
SECRET_<NAME_OF_THE_TARGET_VARIABLE>=vault:<vault_secret_path>:<key_name>

As a result environment variable <NAME_OF_THE_TARGET_VARIABLE> will be set to the content of the secret value

## Usage

```

Usage:
  vault-wrapper [flags]

Flags:
  -a, --auth_path string      login url for vault kerberos plugin
  -h, --help                  help for vault-wrapper
  -t, --keytab string         Path to keytab file
  -k, --krb5 string           Path to krb5.conf (default "/etc/krb5.conf")
  -r, --realm string          kerberos realm for user
  -p, --secret_path strings   Secret path in Vault path/to/secret:key
  -s, --spn string            Vault Service Principal Name
  -o, --token-only            Only output vault token (no execution wrapped)
  -u, --user string           user login
  -v, --vault_addr string     Vault address

Example:
  ./vault-wrapper \
    --vault_addr=http://localhost:8200 \
    --auth_path=/v1/auth/kerberos/login \
    --user=vault \
    --realm=KERBEROS.REALM \
    --spn=HTTP/localhost:8200 \
    --keytab=vault.keytab \
    "command arg arg arg"

```

## Example

```
$ export SECRET_VAULT_SECRET1=vault:keys/boss:private
$ export SECRET_OTHER_SECRET2=vault:accounts/root:password
$ ./vault-wrapper --vault_addr=http://localhost:8200 --auth_path=/v1/auth/kerberos/login --user=vault --realm=KERBEROS.REALM --spn=HTTP/localhost:8200 --keytab=vault.keytab "echo 'something'"
2020/05/29 11:00:33 [VaultWrapper] Authenticating to Vault...
GOKRB5 Client: 2020/05/29 11:00:33 TGT session added for KERBEROS.REALM (EndTime: 2020-05-29 19:00:33 +0000 UTC)
GOKRB5 Client: 2020/05/29 11:00:33 using SPN HTTP/localhost:8200
GOKRB5 Client: 2020/05/29 11:00:33 ticket added to cache for HTTP/localhost:8200 (EndTime: 2020-05-29 19:00:33 +0000 UTC)
2020/05/29 11:00:33 [VaultWrapper] Collecting secrets...
2020/05/29 11:00:33 [VaultWrapper] Launching command:  [sh -c echo 'something']
2020/05/29 11:00:33 [VaultWrapper] Running with PID:  20747
...<Output of the command wrapped>
2020/05/29 11:00:33 [VaultWrapper] Process exited. Exit code:  0
$ env
...
SECRET_OTHER_SECRET2=vault:accounts/root:password
SECRET_VAULT_SECRET1=vault:keys/boss:private
OTHER_SECRET2=pangolin
VAULT_SECRET1=bat
...
```
