# service-mesh-hub 

For operators, the Service Mesh Hub provides a dashboard see the state of the service 
meshes on your cluster, and to discover, deploy, and manage extensions for those meshes. 

For developers, the Service Mesh Hub provides the ability to write and distribute new 
service mesh extensions, solving problems related to enforcing dependency requirements 
and customizing installation manifests based on what meshes have been deployed to the 
cluster. 

## Repo Structure
```
├── api
│   └── v1
├── extensions
│   └── v1
│       ├── flagger
│       ├── gloo
│       ├── glooshot
│       └── kiali
└── pkg/
    ├── kustomize
    └── render

```

#### Api

The api folder contains the API definitions for the service mesh hub resource CRDs. 
These resources are represented as `.proto` files for ease of use and understanding.
The corresponding `.go` files contain generated go representations of the protobuf objects.

#### Extensions

The extensions folder is the main extension registry and the service-mesh-hub. What this means is that
by default the service-mesh-hub will use this folder to search for available/installable meshes.
This folder is where developers will define their application specs which tell the service-mesh-hub
operator how to install the application. For more information on this spec see [here](api/v1/registry.proto). The v1 in this case
corresponds to the version of the API. As the service-mesh-hub progresses the APIs may change, so to
ensure consistency within a given API version, all of the application specs for that version will be kept
in it's corresponding extensions folder.

To find out more about the structure of a mesh extension folder such as `gloo` or `flagger` above
see our documentation [here](api/v1/registry.proto)

#### pkg

Similar to other go projects, pkg houses all of the go code used within the project itself. Contained within
are multiple utilities and libraries pertaining to the rendering and management of kubernetes resources 
and manifests.

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
    --type EXTENSION \
    --print-manifest
```

```bash
GITHUB_TOKEN=`cat ~/github/token/file` go run main.go validate \
    --name kiali \
    --flavor istio \
    --version 0.12 \
    --type EXTENSION \
    --print-manifest
```
