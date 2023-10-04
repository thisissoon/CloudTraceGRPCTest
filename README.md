### Cloud Trace gRPC demo

A simple gRPC server using the [helloworld example](https://github.com/grpc/grpc-go/tree/master/examples/helloworld) to demonstrate instrumenting a grpc server with Open Telemetry and Google Cloud Trace. 

The `SayHello` will automatically fail with an internal grpc error to demonstrate the status handling in Cloud Trace.

Prerequisites:
* A GCP project with tracing enabled
* A GCP user logged in, with the correct IAM permission to create/view traces

To run it:
1. PROJECT_ID=<gcp_project_id> go run *.go
2. See that the call errored with this log message `could not greet: rpc error: code = Internal desc = force error`
3. Go to the Cloud Trace dashboard and observe that span appears normal, it's not coloured red.
