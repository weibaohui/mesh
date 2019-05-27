docker save istio/sidecar_injector:1.1.7 > sidecar_injector.tar
docker save istio/proxyv2:1.1.7 > proxyv2.tar
docker save istio/proxy_init:1.1.7 > proxy_init.tar
docker save istio/pilot:1.1.7 > pilot.tar
docker save istio/mixer:1.1.7 > mixer.tar
docker save istio/kubectl:1.1.7 > kubectl.tar
docker save istio/galley:1.1.7 > galley.tar
docker save istio/citadel:1.1.7 > citadel.tar

docker save  istio/examples-bookinfo-reviews-v3:1.13.0 > examples-bookinfo-reviews-v3.tar
docker save  istio/examples-bookinfo-reviews-v2:1.13.0 > examples-bookinfo-reviews-v2.tar
docker save  istio/examples-bookinfo-reviews-v1:1.13.0 > examples-bookinfo-reviews-v1.tar
docker save  istio/examples-bookinfo-details-v1:1.13.0 > examples-bookinfo-details-v1.tar
docker save  istio/examples-bookinfo-ratings-v1:1.13.0 > examples-bookinfo-ratings-v1.tar
docker save  istio/examples-bookinfo-productpage-v1:1.13.0 > examples-bookinfo-productpage-v1.tar
