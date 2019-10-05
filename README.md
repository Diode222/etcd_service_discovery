## etcd_service_discovery
Service discovery based on etcd.

A example has been implemented in main.go, it will start a service, and this service will be registered in etcd. 
client_main() in main.go is actually the impletation of client, comment the main function and rename client_main as main, 
you can start a client which can request server's pb service.