# Controller to Update the LoadBalancer Status of a Service

This controller updates the `Service.Status.Status.LoadBalancer` field based on the Inputs passed to the controller.

The controller accepts inputs through ENV varibles and below are the explanations for the ENV variables.

| ENV Variables          | Explanation           |
| :-------------: |:-------------:|
| `SERVICE_NAME` | Name of the Service to watch for |
| `SERVICE_NAMESPACE` | Namespace of the Service |
| `EXTERNAL_ENDPOINT_TYPE_ANNOTATION` | Annotation to look for to get the Endpoint type. Example: `status.service.com/endpoint-type` |
| `ENDPOINT_ANNOTATION` | Annotation to look for to get the Endpoint. Example: `status.service.com/endpoint` |

