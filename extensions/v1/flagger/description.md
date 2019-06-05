Flagger is a Kubernetes operator that automates the promotion of canary deployments using Istio or App Mesh routing
for traffic shifting and Prometheus metrics for canary analysis. The canary analysis can be extended with webhooks
for running system integration/acceptance tests, load tests, or any other custom validation.

Flagger implements a control loop that gradually shifts traffic to the canary while measuring key performance
indicators like HTTP requests success rate, requests average duration and pods health. Based on analysis of the KPIs
a canary is promoted or aborted, and the analysis result is published to Slack.