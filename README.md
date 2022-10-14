# ClusterIssuer for Regru API

This solver can be used when you want to use cert-manager with Regru API. API documentation is [here](https://www.reg.ru/reseller/api2doc).


## Install cert-manager (optional step)
If you need install cert-manager in your kubernetes cluster, you can use [command](https://cert-manager.io/docs/installation/) from official documentation.

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

## Install Webhook
```shell
git clone https://github.com/flant/clusterissuer-regru.git
```

You must edit file `values.yaml` in repository by specifying the fields `zone`, `image`, `user`, `password`, for example:
```yaml
issuer:
  zone: my-domain-test.ru
  image: ghcr.io/flant/cluster-issuer-regru:1.0.0
  user: my_user@example.com
  password: my_password
```
where `user` and `password` - are credentials for an authentication of the REG.RU

You must complete commands for install webhook.

```shell
cd clusterissuer-regru
helm install -n cert-manager regru-webhook ./helm
```

## Create ClusterIssuer

Create the file  `ClusterIssuer.yaml` with the contents:
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
and run command

```shell
kubectl create -f ClusterIssuer.yaml
```

### Credentials
In order to access the HTTP API, the webhook needs `user` and `password`.

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

## Create a certificate

Create the file `certificate.yaml` with contetns:

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
    -  *.my-domain-test.ru
```