---
title: Troubleshooting
menuTitle: Troubleshooting
weight: 1
---

Basic troubleshooting tips for the Service Mesh Hub.

## Extension Page

##### No Available Extensions
If no extensions are available, or the extensions page throws an error there are a few possible solutions.


* Add a github token which has access to the given registry

	* Refer to the doc [here](../../installation/authorization) on adding a github token

* bounce the api-server pod

	* Run the following command
```bash
kubectl delete pod -n <hub-namespace> smm-apiserver-<uniqueid of your particular pod>
```
