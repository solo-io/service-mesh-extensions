A Microservice Architecture breaks up the monolith into many smaller pieces that are composed together. Patterns to secure the communication between services like fault tolerance (via timeout, retry, circuit breaking, etc.) have come up as well as distributed tracing to be able to see where calls are going.

A service mesh can now provide these services on a platform level and frees the application writers from those tasks. Routing decisions are done at the mesh level.

Kiali works with Istio to visualize the service mesh topology, features like circuit breakers or request rates.

Kiali also includes Jaeger Tracing to provide distributed tracing out of the box.

To try, run the following command:
```
kubectl port-forward --namespace <install namespace> deployment/kiali 20001:20001
```
then open http://localhost:20001 and authenticate with username `admin` and password `admin`.