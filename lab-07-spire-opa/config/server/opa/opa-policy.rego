package envoy.authz

import input.attributes.request.http as http_request
import input.attributes.source.principal as source_svid

default allow := false

allow {
    http_request.method == "GET"
    source_svid == "spiffe://coastal-containers.example/app/manifest/pilot-boat-0"
}