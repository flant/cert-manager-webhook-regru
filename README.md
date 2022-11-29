# ClusterIssuer for the Regru API

### Motivation

cert-manager automates the management and issuance of TLS certificates in Kubernetes clusters. It ensures that certificates are valid and updates them when necessary.

A certificate authority resource, such as ClusterIssuer, must be declared in the cluster to start the certificate issuance procedure. It is used to generate signed certificates by honoring certificate signing requests.

For some DNS providers, there are no predefined CusterIssuer resources. Fortunately, cert-manager allows you to write your own ClusterIssuer.

This solver allows you to use cert-manager with the Regru API. Documentation on the Regru API is available [here](https://www.reg.ru/reseller/api2doc).

# Usage

### Preparation

You must check access to the Regru API from your IP(s). You should use this command:

```shell
curl "https://api.reg.ru/api/regru2/zone/get_resource_records?input_data=%7B%22username%22%3A%22USER_NAME%22%2C%22password%22%3A%22PASSWORD_STRING%22%2C%22domains%22%3A%5B%7B%22dname%22%3A%22ZONE_NAME%22%7D%5D%2C%22output_content_type%22%3A%22plain%22%7D&input_format=json"

```
where `USER_NAME` and `PASSWORD_STRING` are your credentials to access the Regru API, and `ZONE_NAME` is your domain.

If you doesn't have access, you get an error:

```
{
   "charset" : "utf-8",
   "error_code" : "ACCESS_DENIED_FROM_IP",
   "error_params" : {
      "command_name" : "zone/get_resource_records"
   },
   "error_text" : "Access to API from this IP denied",
   "messagestore" : null,
   "result" : "error"
}
```
You need to add your IP(s) to the whitelist in your personal account.

### Install cert-manager (*optional step*)

**ATTENTION!** You should not delete the cert-manager if you are already using it.


Use the following command from the [official documentation](https://cert-manager.io/docs/installation/) to install cert-manager in your Kubernetes cluster:

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/VERSION/cert-manager.yaml
```
*  where `VERSION` is necessary version (for example, v1.10.1 )

### Install the webhook
```shell
git clone https://github.com/flant/clusterissuer-regru.git
```

Edit the `values.yaml` file in the cloned repository and enter the appropriate values in the fields `zone`, `image`, `user`, `password`. Example:
```yaml
issuer:
  zone: my-domain-test.ru
  image: ghcr.io/flant/cluster-issuer-regru:1.0.1
  user: my_user@example.com
  password: my_password
```

You must also specify your namespace with the `cert-manager`.

```yaml
certManager:
  namespace: my-namespace-cert-manager
  serviceAccountName: cert-manager
```

Here, `user` and `password` are credentials you use to authenticate with REG.RU.

Next, run the following commands for the install webhook.

```shell
cd clusterissuer-regru
helm install -n cert-manager regru-webhook ./helm
```

### Create a ClusterIssuer

Create the `ClusterIssuer.yaml` file with the following contents:
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
and create the resource:

```shell
kubectl create -f ClusterIssuer.yaml
```

#### Credentials

You have to provide a `user` and `password` for the webhook so that it can access the HTTP API.

Note that we use `regru-password` as the secret reference name in the `ClusterIssuer` example above. If you use a different name for the secret, make sure to edit the value of `regruPasswordSecretRef.name`.

The secret for the above example would be as follows:
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

Create the `certificate.yaml` file with the following contents:

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

# Community

Please feel free to contact us if you have any questions.

You're also welcome to follow [@flant_com](https://twitter.com/flant_com) to stay informed about all our Open Source initiatives.

# License

Apache License 2.0, see [LICENSE](LICENSE).
