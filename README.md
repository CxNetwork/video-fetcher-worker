# video-fetcher-worker (CxNetwork Open-Source)
This is a microservice, part of the Cx Network backend. It runs on just 1 pod on our Kubernetes cluster.

This microservice is responsible for fetching the latest videos of streamers in the network and updating the database with them.

video-fetcher-worker is designed to run on Google Cloud Kubernetes Engine.

# Running
Build:  
`docker build -t gcr.io/[PROJECT ID]/videofetcher .`

Push image to Google Container Registry:  
`gcloud docker -- push gcr.io/[PROJECT ID]/videofetcher`

Attach to Kubectl:  
`gcloud container clusters get-credentials [CLUSTER NAME] --zone [ZONE] --project [PROJECT ID]`

Create the deployment:  
`kubectl create -f video-fetcher-backend.yaml`

You should never need to scale the deployment to more than 1 pods.