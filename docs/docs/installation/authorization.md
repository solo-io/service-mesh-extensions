---
title: Authentication with Github
menuTitle: Authentication with Github
weight: 2
---

Many resources in the Service Mesh Hub are stored in Github. As you install or modify service mesh extensions
you may encounter rate limiting issues.

You can resolve this by creating a Kubernetes secret with a valid Github token.

### Get a Github token

If you do not already have a Github token, you can acquire one as described in the [Github docs](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)

### Create a Kubernetes secret

With your token, run the following command to create a secret:
```bash
GITHUB_TOKEN=<token-value> ./install/create-secret.sh
```

Alternatively, you can call this command:
```bash
kubectl create secret generic github-token \
    -n sm-marketplace \
    --from-literal=token=<token-value>
```

### Restart the Service Mesh Hub pods

Now that you have created the necessary secret, you need to restart the Service Mesh Hub pods so they can
use it.

If you have not installed anything yet, you can restart all the pods in the sm-marketplace namespace.
```bash
kubectl delete pod -n sm-marketplace --all
```

Alternatively, just delete the `smm-operator` and `smm-apiserver` pods:
```bash
kubectl get pods -n sm-marketplace
kubectl delete pod -n sm-marketplace smm-operator-<ABC> # substitute actual pod name
kubectl delete pod -n sm-marketplace smm-apiserver-<DEF> # substitute actual pod name
```

After these pods have been restarted you should have no further issues with Github ratelimits.
