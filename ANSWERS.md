### 1. How would you measure the performance of your service?

I've enjoyed using [Datadog][dd] for monitoring application metrics.

asciist would use Datadog's [Go integration][dd-go] to publish a few metrics:

- number of requests processed, broken down by response status
- time spent processing requests
  - overall time, broken down by HTTP status
  - time spent decoding the JSON payloads
  - time spent in [image.Decode][decode], broken down by image format
  - time spent in ASCII art conversion of decoded images
- number of bytes processed of input images
- number of pixels processed of input images

This would allow an overall view of application health, including:

- average, median, and 90th/95th/99th percentile times spent waiting for your ASCII art
- percentage of requests whose response takes longer than a threshold time considered “too long to wait”
- processing time per pixel, facilitating the detection of performance improvements/regressions and helping to understand aberrations (e.g., a DoS attack using specially-crafted images exploiting an edge case in image decoding)

[dd]: https://www.datadoghq.com/
[dd-go]: https://godoc.org/github.com/DataDog/datadog-go/statsd
[decode]: https://golang.org/pkg/image/#Decode

### 2. What are some strategies you would employ to make your service more scalable?

Fortunately, this is a stateless service that can be readily horizontally scaled. It is bounded only by available CPU time and memory.

To take best advantage of that, I would split asciist into two services: a **web tier** that would be responsible for accepting client requests and serving responses, and a **worker tier** that would be responsible for image processing.

We could then readily take advantage of AWS managed services to let the app be scaled in response to demand without our intervention, namely autoscaling, SQS, and S3.

#### Architecture overview

Upon receiving a request, the web tier will:

- [validate][image-config] that the received image can be processed, and send back a 400 if it can't
- compute an SHA2 hash of the image's contents
- check if there's already an image matching that hash in S3
  - if not, upload the image to S3
  - if there is, try to get a result from S3 matching the image hash and desired ASCII art width
    - if that worked, send it to the client immediately
- publish a message to an SQS queue, including:
  - the S3 object key of the image
  - the user-desired width of the ASCII art
  - an deadline; the time at which the web tier will consider the request to have timed out
- wait for a worker to indicate that the image has been processed
- fetch the ASCII art from S3 and serve it to the client

Upon receiving a message from SQS, the worker tier will:

- check the deadline; if the web tier already gave up on this request, delete the message from SQS and return immediately
- fetch the image from S3
- convert it to ASCII art
- put the art in S3 at `/{image hash}/width/{width}`
- notify the web tier that the request has been processed
- delete the message from SQS

[image-config]: https://golang.org/pkg/image/#Config

#### Notifying the web tier

One question not addressed above is: how do workers let web instances know when a request they're waiting for has finished?

##### Direct approach

Since each response would map back to one specific Web process, an easy way would be to include a callback URL in the SQS messages giving a direct host:port for the web process.

Web processes would start a shared goroutine to maintain a local map of *SQS message ID* -> *result channel*. The HTTP request handlers would use a channel to talk to this goroutine to inform it of new and completed worker requests. When waiting for an image, the request handler would `select` across the result channel and a timeout.

##### SQS approach

If the application is being deployed in containers, exposing a direct route to the Web process may be complicated. In that case, each process could create its own SQS queue on startup, and send the queue URL in messages to workers.

#### Notes

- Both web and worker tier instances would be deployed in auto-scaling groups
  - Worker instances can be scaled based on the SQS request queue size
  - Web instances can be scaled based on EC2 metrics (CPU, memory)
- An S3 lifecycle rule can be used to delete images and results automatically

### 7. If you wanted to migrate your scaled-out solution to another cloud provider (with comparable offerings but different API’s) how would you envision this happening? How would deal with data consistency during the transition and rollbacks in the event of failures?

Again, fortuately this is a stateless service: its sole request handler is a function of its inputs.

To transition to a new provider smoothly and with the ability to roll back, I would:

1. Route traffic to the existing provider using a third-party CDN (like Cloudflare).
2. Add the new cloud to the CDN's set of origin servers.
3. Monitor.
  - Problems with the new cloud? Take it out of the CDN, fix, repeat.
4. Remove the old cloud from the CDN and spin down its resources.
5. Optionally update DNS to point to the new cloud natively and spin down the CDN.

If this service did have persistent data, I would need to do a lot of research to feel confident making that transition. Hopefully you're using not using proprietary services (e.g. DynamoDB) that have no direct analogue on the new provider.

But, conceptually, I'd think it's similar to the above with additional steps:

1. Spin up data store(s) in the new cloud, and make them replication secondaries of the primaries in the original cloud.
2. Put the service in the original cloud briefly into a read-only mode and wait for the new cloud to process the latest transactions.
3. Switch roles: promote databases in the new cloud to primary, and have databases in the old cloud replicate from them.
4. Switch Cloudflare to point to the new cloud.
5. Bring the service out of read-only.
6. Monitor.
  - Problems with the new cloud? Invert what you just did.
7. Eventually wind down the old cloud. Maybe keep secondary databases there for disaster recovery purposes.
