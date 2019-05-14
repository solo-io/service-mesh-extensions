# service-mesh-hub 

For operators, the Service Mesh Hub provides a dashboard see the state of the service 
meshes on your cluster, and to discover, deploy, and manage extensions for those meshes. 

For developers, the Service Mesh Hub provides the ability to write and distribute new 
service mesh extensions, solving problems related to enforcing dependency requirements 
and customizing installation manifests based on what meshes have been deployed to the 
cluster. 


# Validate your extension before deploying


If you are creating a new extension or modifying an existing one, you can verify that your specification
will be accepted prior to submission using the following command line tool.

Just specify the `--name`, `--flavor`, and `--version`.

If you want to see a preview of the corresponding chart, pass the `--print-manifest` flag.

Here are two examples:

```bash
GITHUB_TOKEN=`cat ~/github/token/file` go run main.go validate \
    --name glooshot \
    --flavor istio \
    --version 0.0.2 \
    --print-manifest
```

```bash
GITHUB_TOKEN=`cat ~/github/token/file` go run main.go validate \
    --name kiali \
    --flavor istio \
    --version 0.12 \
    --print-manifest
```
