# Lab 6: Navigating Basic Authorization with Open Policy Agent

## Prerequisites

- A 64-bit Linux environment (AMD, Intel, or Arm).
- Basic familiarity with Docker containers and `docker` commands.
- [jq](https://jqlang.github.io/jq/) installed on your device.

## Introduction

Coastal Containers Ltd. is piloting a new system to manage access to their `port-records` database, which contains vital information about vessel manifests. Fleet Supervisors and Ship Captains alike need to access these manifests for operational and logistical purposes. However, Coastal Containers wants to ensure that the captains can only access records relevant to their assigned vessels, while fleet supervisors, due to their managerial role, should have unrestricted access to all records.

### Open Policy Agent

In order to implement the aforementioned business logic, Coastal Containers Ltd. has opted to explore the usage of [Open Policy Agent](https://github.com/open-policy-agent/opa), an open source, general purpose policy engine. Open Policy Agent (aka OPA) is a tool which can make policy decisions on structured data using the [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) policy language. Throughout the course of this hands-on exercise, we will be using Rego to implement authorization through Open Policy Agent.

- **Design Architecture:** To better understand how OPA works, you can check out the [OPA Overview](https://www.openpolicyagent.org/docs/latest/#overview), which outlines the design architecture for Open Policy Agent and provides a simple [example](https://www.openpolicyagent.org/docs/latest/#example) to demonstrate how it all works.

- **Rego:** As mentioned previously, Rego is the native policy language used by OPA, which allows for convenient data querying and policy definitions. Rego is built to be declarative, allowing authors to focus on outcome as opposed to execution. Before diving into this exercise, it is reccommended to check out the [Policy Language](https://www.openpolicyagent.org/docs/latest/policy-language/) and [Policy Reference](https://www.openpolicyagent.org/docs/latest/policy-reference/) documentation to build a baseline understanding of how Rego works. On top of these resources, OPA offers the [OPA Playground](http://play.openpolicyagent.org/), an interactive environment where you can explore packaged examples, test your own Rego policies, and validate their output.

### Business Logic

The logic behind policy decisions is fairly straightforward for this scenario:

- **Fleet Supervisors:** Should have unrestricted access to view all vessel records in the `port-records` database.

- **Ship Captains:** Should only access records of vessels assigned to them, ensuring they cannot view or modify records irrelevant to their operational duties.

*📝Note: For the purposes of our demonstration, we will abstract the details of where and how policy decisions are enforced.*

In a more complex example, policy decisions are often enforced by [Policy Enforcement Points](https://csrc.nist.gov/glossary/term/policy_enforcement_point) (PEP), which could be something like a network proxy filtering traffic before it hits the application (e.g., [Envoy Proxy](https://www.openpolicyagent.org/docs/latest/envoy-introduction/)). However, in our use-case example, all you need to know is that HTTP requests are made to the 'example' application, and contain the following:

- a `bearer` field containing a [JWT](https://jwt.io/introduction) (and associated claims) about the authenticated calling user. 

- a `vessel_id` field which represents the Coastal Containers vessel record being accessed.

- an `action` to be performed on the record, such as `read`, `update`, or `delete`. For the purposes of this scenario, we won't go into depth about the different actions, and will stick with `read` as a basis for testing policy decisions.

The `port_data` stored within the `port-records` database is made available to OPA as an [input](https://www.openpolicyagent.org/docs/latest/envoy-primer/#input-document), which we will explore in-depth later in this lab.

### Preparing Your Environment

Before you cast off, prepare your ships to sail by setting your working directory in [lab-06-opa-basics](../lab-06-opa-basics/) as an environment variable:

```bash
export LAB_DIR=$(pwd)
```

This will make issuing commands easier in the following steps of this exercise, and will reduce the possibility of reference errors.

## Step-by-Step Instructions

### Step 1: Obtain an Example JWT

To simulate user authentication in our simplified scenario, we will need to create example JWTs (JSON Web Tokens). In order to do so, we will invoke the inbuilt `io.jwt.encode_sign` function in Rego. Explore how this is done within the [create_jwt.rego](./create_jwt.rego) file, and utilize the [Policy Reference](https://www.openpolicyagent.org/docs/latest/policy-reference/) as necessary.

#### create_jwt.rego Breakdown

Before we can create the example JWTs, we must understand how the [create_jwt.rego](./create_jwt.rego) works with Rego. At a baseline, this policy file works to generate JWTs with pre-specified claims that indicate the `fleet_supervisor` or `ship_captain` role(s), and other relevant metadata. At the beginning, the [create_jwt.rego](./create_jwt.rego) file is initialized with a hierarchical `package` name, ensuring that Rego policies and rules are organized based on their functionality. After this, the built-in functions [ceil] and [time.now_ns], are used to define an `expiry_time` for our JWTs. This `expiry_time` is set to one day from the point of creation, however, this should **NOT** be used in production scenarios as JWTs are intended to be short-lived.

Throughout this file, you will notice the usage of equality operators such as `:=`, `==`, & `=`. Natively, Rego supports three kinds of equality:

- Assignment (`:=`) is used to assign values to variables. Assigned variables are locally scoped to the rule that they are set within and 'shadow' global variables. An in-depth explanation of this works can be found [here](https://www.openpolicyagent.org/docs/latest/policy-language/#assignment-) on the official docs.

- Comparison (`==`) is used to check if two values are equal within a rule, this is recursive and semantic. An in-depth explanation of how this works can be found [here](https://www.openpolicyagent.org/docs/latest/policy-language/#comparison-) on the official docs.

- Unification (`=`) is used to combine assignment and comparison. Rego will assign one or more variables to make the defined comparison true, effectively letting you query values for variables that make an expression true. An in-depth explanation of how this works can be found [here](https://www.openpolicyagent.org/docs/latest/policy-language/#unification-) on the official docs.

Now, with the `expiry_time` set for our example tokens, we can now create JWTS for the Fleet Supervisor and Ship Captain using the invoked `io.jwt.encode_sign` function. Keep in mind that we have hardcoded an RSA key pair in the [create_jwt.rego](./create_jwt.rego) file, this is **NOT** safe for production and should only be used for demonstration purposes as the RSA private key can be easily stolen.

To obtain the Fleet Supervisor's JWT (`spv_token`), try using the [opa eval](https://www.openpolicyagent.org/docs/latest/cli/#opa-eval) command by running OPA in a Docker container:

```shell
docker run -v ${LAB_DIR}:/example openpolicyagent/opa:0.57.0 eval -d /example/create_jwt.rego 'data.example.jwt.spv_token' | jq '.result[0].expressions[0].value'
```

Functionally, this command loads the [create_jwt.rego](./create_jwt.rego) into OPA using the `-d` or `--data` flag, and queries the value of `data.example.jwt.spv_token`. This works through the [OPA Document Model](https://www.openpolicyagent.org/docs/latest/philosophy/#the-opa-document-model), demonstrating how the hierarchical `package` and associated rules sit under OPA's `data` document.

The expected output should be:

```log
"eyJhbGciOiAiUlMyNTYifQ.eyJhdWQiOiAicG9ydC1yZWNvcmRzIiwgImV4cCI6IDE2OTcwMzc1NDksICJpc19zdXBlcnZpc29yIjogdHJ1ZSwgImlzcyI6ICJodHRwczovL2lkcC5jb2FzdGFsLWNvbnRhaW5lcnMuZXhhbXBsZSIsICJzdWIiOiAiZmxlZXRfc3VwZXJ2aXNvciJ9.OKq69mj0Z22l4I2pWKSr4xaErVKBEQcdOaCYUi3sckUmjixYFb4nZGRXFp2eSPlYdhDiqldgBkrE9W1--8Soluemamg1WHd4jrPtKwHKwFHPAkrH4TUTHJ-3wXIeWr8WRXDiulvBOAd2w4Wmq0fMUo3iwTnN5M67dBUmtqSX03tkwnL7QdIHUwpTGYaBm79N5RiOo_vw7HPtkZv6nLTd0LYT9jui_EpL4l-jHQxlp8omuI9FupjHkA1tWRtEh3ny_prSgntV1X277_EkWmJh0TrORQDoZ390gxaDSTcvfxxIdICsdogG_UT4mBPhqalAByUigVgTmgngykm4qSsfxw"
```

To verify the token includes the defined claims, copy and paste the outputted JWT into the `Encoded` field of [jwt.io](https://jwt.io/). This operation will show the subsequently decoded fields or 'claims', such as the `is_supervisor` claim which denotes if a role is a Fleet Supervisor per our scenario's business logic. For the Ship Captain business logic, however, we will need to better understand how OPA manages external data.

### Step 2: Inspect the Data Structure

In the previous step, we have walked through a simple demonstration of how OPA can handle external data via JWT tokens. However, per our scenario of the `port-records` database which represents a *relatively* static store of vessel information, which should change infrequently, and can be reasonably stored in-memory all at once, we can replicate it in bulk via OPA's [bundle](https://www.openpolicyagent.org/docs/latest/external-data/#option-3-bundle-api) feature. Through this approach, both policies and external data can be added to a bundle (`tar.gz` file). OPA can then consume the packaged bundle via a bundle server.

Reference our example [bundle](./bundle/) directory to see where our external data sits within OPA's `data` document. Within our setup, the [data.json](./bundle/port_data/data.json) file represents the `port-records` database, and is located within the [port_data](./bundle/port_data/) directory. This means that specific role data can be accessed at `data.port_data.roles`, or `data.port_data.vessels` for vessel-related queries.

### Step 3: Inspect the Policy File

Keeping the data structure and example JWTs in mind, we are now in a position to write and define policies implementing Coastal Containers business logic. First, inspect the [policy.rego](./bundle/example/authz/policy.rego) file to see how we intend to do this. Notice that the outcome of our policy decision(s) will be encapsulated in the value of `allow`. Due to the hierarchical nature of packages and the OPA `data` document, this information becomes available at `data.example.authz.allow`.

For Coastal Containers Ltd., the input provied to OPA will adhere to the following structure:

```shell
{ "input": { "bearer": "<JWT>", "action": "read", "vessel_id": 1 } }
```

Within our policy, we can refer to `input.bearer`, `input.action`, and `input.employee_id`. However, adding an `import` statement at the top of the policy file (e.g., `import input.bearer`), means we can refer to `bearer` in the encapsulated Rego rules.

Fundamentally, OPA policies are formed as a collection of rules, where rules take the form of `assignment if { condition(s) }`. An example of this within our policy is:

```rego
allow {
    token_is_valid
    role_is_supervisor
}
```

The rule body between `{}` is a collection of assignments and expressions. `allow` will evaluate to true if a logical AND of all the assignments and expressions is true. If an assignment is false or undefined, allow is also undefined. As such, we need to set `default allow := false` in the policy, so that `allow` can only be true if one of the rules evaluates to true - otherwise it will be false, but never undefined. In this way, multiple `allow` rules represent a logical OR.

You are encouraged to read through and understand the rest of the policy, referring out to the [OPA documentation](https://www.openpolicyagent.org/docs/latest/) if necessary. 

### Step 4: Build the Policy and Data Bundle

With a baseline understanding of how the underlying tooling works, and how we plan to implement it, we can start this demo by building the policy and data bundle.

To do so, run the following docker command:

```shell
docker run -v ${LAB_DIR}:/example openpolicyagent/opa:0.57.0 build \
    --bundle /example/bundle \
    -o example/bundle.tar.gz
```

Once this is completed, you should see the packaged `bundle.tar.gz` file within the root [lab-06-opa-basics](../lab-06-opa-basics/) directory.

### Step 5: Run OPA in a Docker Container

Now, to get OPA up and running locally, we will run OPA as a server in a Docker container:

```shell
docker run --name opa-server --rm -p 8181:8181 -d -v ${LAB_DIR}:/example \
    openpolicyagent/opa:0.57.0 run --server \
    --bundle /example/bundle.tar.gz \
    --addr 0.0.0.0:8181
```

Once executed, you should see an outputted container ID like this:

```log
f3b4d5da79909e1f577a34dfab869517c58a8814e88c10274b2bc7e6bb572bb1
```

This indicates that the server is running as a docker container (in the background via the `-d` detach flag), and listening on `http://0.0.0.0:8181`. To view the logs of the newly created `opa-server` container, run:

```shell
docker logs opa-server
```

You should see the output:

```log
{"addrs":["0.0.0.0:8181"],"diagnostic-addrs":[],"level":"info","msg":"Initializing server.","time":"2023-10-10T15:10:57Z"}
```

After the `opa-server` container begins waiting for requests, we can move onto the next steps.

*📝Note: To view the list of running Docker containers on your device, you can run:*

```shell
docker ps
```

### Step 6: Export JWTs as Environment Variables

Navigate to the root [lab-06-opa-basics](../lab-06-opa-basics/) directory by running:

```shell
cd $LAB_DIR
```

Once here, we will export our Fleet Supervisor & Ship Captain JWTs as environment variables for the sake of convenience and reusability.

To do this for the Fleet Supervisor token (`spv_jwt`), run:

```shell
export SPV_JWT=$(docker run -v ${LAB_DIR}:/example openpolicyagent/opa:0.57.0 eval -d /example/create_jwt.rego 'data.example.jwt.spv_token' | jq '.result[0].expressions[0].value')
```

To do this for the Ship Captain token (`cpt_jwt`), run:

```shell
export CPT_JWT=$(docker run -v ${LAB_DIR}:/example openpolicyagent/opa:0.57.0 eval -d /example/create_jwt.rego 'data.example.jwt.cpt_token' | jq '.result[0].expressions[0].value')
```

Functionally, these commands are loading the [create_jwt.rego](./create_jwt.rego) file into OPA via the `-d` or `--data` flag, and querying the tokens stored at `data.example.jwt.spv_token` & `data.example.jwt.cpt_token`. As mentioned previously, this query utilizes the hierarchical nature of OPA packages (specifically the `example.jwt` package).   

### Step 7: Test the Policy Decisions

In order to evaluate the policy decisions per our business logic, we can make POST requests to the OPA server in the following format:

```shell
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$CPT_JWT"', "action": "read", "vessel_id": 1}}' 0.0.0.0:8181/v1/data/example/authz/allow
```

*📝Note: Through the form of this POST request, we are providing the relevant input via the parameters in the body of the request, and we are querying the value of `data.example.authz.allow` via OPA's [Data API](https://www.openpolicyagent.org/docs/latest/rest-api/#data-api).*

Now, to apply this request format and validate our scenario's business logic, try answering the following questions by issuing the associated POST requests. 

Can the Fleet Supervisor view the `maritime-mover` vessel?

```shell
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$SPV_JWT"', "action": "read", "vessel_id": 1}}' 0.0.0.0:8181/v1/data/example/authz/allow
```

To further test our business logic, try this command again and run it for `"vessel_id": 2` & `"vessel_id": 3`, ensuring the Fleet Supervisor can view all of them. 

The expected output is:

```log
{"result":true}
```

Next, can the Ship Captain view it's assigned vessel (`maritime-mover`)?

```shell
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$CPT_JWT"', "action": "read", "vessel_id": 1}}' 0.0.0.0:8181/v1/data/example/authz/allow
```

The expected output is:

```log
{"result":true}
```

Finally, can the Ship Captain view a vessel it is not assigned to?

```shell
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$CPT_JWT"', "action": "read", "vessel_id": 2}}' 0.0.0.0:8181/v1/data/example/authz/allow
```

To further test our business logic, try this command again and run it for `"vessel_id": 3`, ensuring the Ship Captain can only see the assigned `maritime-mover` vessel.

The expected output is:

```log
{"result":false}
```

Keep in mind that, at anytime, you can view the `opa-server` logs by running:

```shell
docker logs -f opa-server
```

This command will 'follow' the `opa-server` Docker container logs via the `-f` or `--follow` flag, and should provide an ouput similar to:

```log
{"addrs":["0.0.0.0:8181"],"diagnostic-addrs":[],"level":"info","msg":"Initializing server.","time":"2023-10-10T20:09:10Z"}
{"client_addr":"172.17.0.1:44962","level":"info","msg":"Received request.","req_id":1,"req_method":"POST","req_path":"/v1/data/example/authz/allow","time":"2023-10-10T20:14:51Z"}
{"client_addr":"172.17.0.1:44962","level":"info","msg":"Sent response.","req_id":1,"req_method":"POST","req_path":"/v1/data/example/authz/allow","resp_bytes":16,"resp_duration":1.279826,"resp_status":200,"time":"2023-10-10T20:14:51Z"}
{"client_addr":"172.17.0.1:44968","level":"info","msg":"Received request.","req_id":2,"req_method":"POST","req_path":"/v1/data/example/authz/allow","time":"2023-10-10T20:14:57Z"}
{"client_addr":"172.17.0.1:44968","level":"info","msg":"Sent response.","req_id":2,"req_method":"POST","req_path":"/v1/data/example/authz/allow","resp_bytes":16,"resp_duration":0.926202,"resp_status":200,"time":"2023-10-10T20:14:57Z"}
{"client_addr":"172.17.0.1:35082","level":"info","msg":"Received request.","req_id":3,"req_method":"POST","req_path":"/v1/data/example/authz/allow","time":"2023-10-10T20:15:10Z"}
{"client_addr":"172.17.0.1:35082","level":"info","msg":"Sent response.","req_id":3,"req_method":"POST","req_path":"/v1/data/example/authz/allow","resp_bytes":16,"resp_duration":1.253941,"resp_status":200,"time":"2023-10-10T20:15:10Z"}
{"client_addr":"172.17.0.1:57136","level":"info","msg":"Received request.","req_id":4,"req_method":"POST","req_path":"/v1/data/example/authz/allow","time":"2023-10-10T20:15:14Z"}
{"client_addr":"172.17.0.1:57136","level":"info","msg":"Sent response.","req_id":4,"req_method":"POST","req_path":"/v1/data/example/authz/allow","resp_bytes":17,"resp_duration":0.916282,"resp_status":200,"time":"2023-10-10T20:15:14Z"}
```

### (Optional) Step 8: Diving Deeper

For those interested in diving deeper into how OPA can integrate with SPIRE and manage policy decisions with Rego, you are encouraged to check out the [envoy-jwt-opa](https://github.com/spiffe/spire-tutorials/tree/main/k8s/envoy-jwt-opa) & [envoy-opa](https://github.com/spiffe/spire-tutorials/tree/main/k8s/envoy-opa) demo exercises located in the [spire-tutorials](https://github.com/spiffe/spire-tutorials/tree/main) repository. Alternatively, you can try your hand at expanding on this demo by building more roles, vessel assignments, and policy decisions. Can you get the policies to work as-expected with your new configurations?

### Step 9: Cleanup

To kill the `opa-server` Docker container, run:

```shell
docker kill opa-server
```

Next, navigate to the root lab directory and remove the `bundle.tar.gz` tarball:

```shell
cd $LAB_DIR && rm bundle.tar.gz
```

## Conclusion

Congratulations aspiring captain! You have helped Coastal Containers Ltd. venture into unexplored waters by running a simple implementation of Open Policy Agent providing policy decisions to their `port-records` datastore. This is a small, but important step towards implementing a robust authorization system into their current shipping systems. In the next lab, we will explore how to bridge turbulent seas and integrate OPA with SPIRE.

To learn more about Open Policy Agent and the Rego policy language, you are highly encouraged to explore the [official documentation](https://www.openpolicyagent.org/docs/latest/), and the existing [OPA Ecosystem](https://www.openpolicyagent.org/ecosystem/) which outlines current integrations with the policy engine.
