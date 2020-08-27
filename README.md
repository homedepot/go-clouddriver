<img src="https://github.com/billiford/go-clouddriver/blob/media/clouddriver.png" width="125" align="left">

# go-clouddriver

go-clouddriver is a rewrite of Spinnaker's [Clouddriver](https://github.com/spinnaker/clouddriver) microservice. It aims to fix severe scaling problems and operational concerns when using Clouddriver at production scale. 

It changes how clouddriver operates by providing an extended API for account onboarding (no more "dynamic accounts") and removing over-complicated strategies such as [Cache All The Stuff](https://github.com/spinnaker/clouddriver/tree/master/cats) in favor of talking directly to APIs.
