# ClusterIssuer for Regru API

This solver can be used when you want to use cert-manager with Regru API. API documentation is [here](https://www.reg.ru/reseller/api2doc).




### ClusterIssuer

Create a `ClusterIssuer` resource as following:
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: regru-dns
spec:
  acme:
    # Email address used for ACME registration. REPLACE THIS WITH YOUR EMAIL!!!
    email: mail@example.com
    # The ACME server URL
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: cert-manager-letsencrypt-private-key
    solvers:
      - dns01:
          webhook:
            config:
              regruPasswordSecretRef:
                name: regru-password
                key: REGRU_PASSWORD
            groupName: {{ .Values.groupName.name }}
            solverName: regru-dns
```

### Credentials
In order to access the HTTP API, the webhook needs an user and a password.

If you choose another name for the secret than `regru-password`, ensure you modify the value of `regruPasswordSecretRef.name` in the `ClusterIssuer`.

The secret for the example above will look like this:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: regru-password
data:
  REGRU_PASSWORD: {{ .Values.issuer.password | b64enc | quote }}
type: Opaque
```

### Create a certificate

Finally you can create certificates, for example:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: changeme
  namespace: changeme
spec:
  secretName: changeme
  issuerRef:
    name: regru-dns
    kind: ClusterIssuer
  dnsNames:
    - example.com
```