# AutoScaler agent

[autoscaler.cloud](https://autoscaler.cloud) is a service that lets you Scale your server fleet automatically depending on the actual load it sustains. Unhealthy servers are also automatically replaced.

This agent sends information about server load to service back-end.
It is automatically installed when a new server is added by the service to an auto scaling pool.
It just push and does not pull any information from the service.
It only supports GNU/Linux systems.
