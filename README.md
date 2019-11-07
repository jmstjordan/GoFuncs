## Exploring Go Cloud Functions and an Authentication Proxy

### Introduction

This repository is meant to be a demonstration of how to set up HTTP Cloud Functions in Google Cloud Platform using Golang and Cloud Firestore. Additionally, it will cover the authentication of these functions using GCP's Extensible Service Proxy. In this example, we will plug Firebase in as our frontend, but as you will see this component is loosely coupled and can be replaced by whatever client you desire.

### First Steps

Unforunately, to do any work in the cloud you need to setup an account and enable billing. If this is your first time using GCP, or have made an account less than a year ago, then you are eligible for $300 of cloud money; so you can enable billing without any fear! That being said, I did not use any of that $300 during this project.

If you want to just explore Go and Cloud Functions, all you will need to test locally is Go version 1.11, and a Cloud Firestore database generated (which you need an GCP account for). When you are ready to setup authentication, that's when you will be going down a big rabbit hole in GCP. So let's start with just some local golang cloud functions and worry about security later!

*In my opinion, GCP's [documentation](https://cloud.google.com/docs/) sets itself apart from other cloud providers. This entire tutorial was done using those docs. When I don't provide links, search these docs for relevant keywords to follow my work.*

### Go and Cloud Functions

Before we begin let me explain the business requirement: I wanted to develop a simple backend allowing users to clock in or clock out. This could be used by a company to track time worked by employees.

1. Once you have an [account](https://cloud.google.com/), go to [GCP Cloud Console](https://console.cloud.google.com), create a project, then enable billing (you'll be fine). Using the search bar at the top to navigate services is very useful.

2. Then, search for Cloud Firestore (again, in the GCP Console) and generate a database.

3. You will need to create two collections, **users** and **clocks**. The **users** is where we still store a document per user, and **clocks** will store a document everytime a user clocks in or out. A collection requires you to have at least one document, so go ahead and create yourself as a user. Then, add an initial clock in for yourself. For the document attributes, follow the templates **User** and **ClockPost** in **get_user.go** and **clockin.go**. Note that each clock in/out will need a reference to a user document.

4. Next, you need to [create a service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts) and give it some sort of elevated permissions. Then, [generate a key](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) tied to that account. This gives us a clean method to authenticate with Firestore without storing credentials in our code. More information on this authentication can be found [here](https://cloud.google.com/functions/docs/securing/authenticating). Download the .json key to your local machine and then set a local environment variable to that location. 

```bash
export GOOGLE_APPLICATION_CREDENTIALS='/Users/<path_to_json>/key.json' 
```

5. Download and install Go, I use version 1.11, but more recent versions should be fine.

6. Now you should have done enough to test the current functionality. We have two tests. The first one retrieves a user from the **users** collection. Change the test code to reflect the user you added in step 3. The second test will actually write a document to **clocks**, so make sure your collections are setup correctly, and the user in **clockin_test.go** reflects a user in the **users** collection. The following command will run both tests:
```bash
go test
```
Or to run them individually:
```bash
go test get_user.go get_user_test.go
```
7. Although there are several ways to deploy cloud functions, I would recomend installing the [GCP SDK](https://cloud.google.com/sdk/install) and deploying them from your local machine. To deploy, run the following commands:
```bash
gcloud functions deploy GetUser --entry-point GetUser --runtime go111 --trigger-http
gcloud functions deploy ClockTime --entry-point ClockIn --runtime go111 --trigger-http
```
The first name (for example ClockTime) will be the name of function in GCP. The entrypoint flag refers to the name of the local function to deploy.

Once you deploy the functions, be sure to address the function's accessibility. Currently by default, all functions are public (although the GCP documentation says this is changing soon). Navigate to your Cloud Functions page, select the function, then go to the **SHOW INFO PANEL**- here you can manage access. I would make these are private until you decide how to handle authorization. You can still trigger these functions for testing in the Cloud Functions page. For demonstration purposes, this is enough until you are ready to actually expose some data and functionality.

### Firebase and Authentication

Now that we have our functions deployed, we need a reasonable way to expose that data without making each of your functions public. There are many [ways to do this](https://cloud.google.com/functions/docs/securing/authenticating), but I chose **Firebase Authentication**. I figured I would need a front-end at some point anyway, so Firebase made the most sense to me based on my use case. One thing to note here is that you can completely bypass these cloud functions and authentication proxy by just using [some of Firebase's basic Firestore interactions](https://firebase.google.com/docs/firestore/quickstart), but as your application scales and you need may more backend logic. And in this case, my goal was to learn some Go and learn about GCP's authentication without Firebase holding my hand.

Setting this process up was pretty involved, and I followed the documentation pretty closely, so I will refer you to the [above link](https://cloud.google.com/functions/docs/securing/authenticating)- navigate down to **Firebase Authentication**.

As a result of authenticating with Firebase, I deployed an Extensinable Service Proxy (ESP) container in GCP's Cloud Run; you can find that tutorial [here](https://cloud.google.com/endpoints/docs/openapi/get-started-cloud-functions), but it is also referenced in the above link. This ESP uses an OpenAPI document which can be found as **openapi-functions.yaml**. This document is essentially the configuration that tells the proxy how to handle authentication, route traffic to our functions, and can also be used as a template to generate a swagger api documentation page.

A few things to callout in this file are the firebase authentication and the modular design this proxy offers. The firebase configuration coupled with the *--cors_allow_origin* flag during deployment allowed me to require all API users to be authenticated through our firebase application. Once setup, the firebase client will simply pass an authentication token to the proxy as part of any request.

Finally, this solution allows us to decouple any backend functionality from the API layer. The API handles routing and authentication, without having to care what downstream functions do. And the functions do the work, and can be changed without caring about how traffic gets to them.

I hope you found this demonstration as useful as I did!
