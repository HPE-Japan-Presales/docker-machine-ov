# Known Issues
## Rancher
### Deploy RKE 1.20.4
Cannot deploy RKE 1.20.4 from this drver normally.  
You can see the etcd issue in log.

```
(error "tls: failed to verify client's certificate: x509: certificate signed by unknown authority (possibly because of "crypto/rsa: verification error" while trying to verify candidate authority certificate "kube-ca")", ServerName "")
```

After restarting etcd container manually, it will be healthy.