# Lab 2: Installing and Configuring SPIRE from Binaries

## Prerequisites

- A AMD, Intel, or Arm 64-bit Linux environment. *Note: it is possible to run within a Docker Container, e.g. `docker run -it -v ${PWD}:/app debian bash` on a Mac.*
- The [OpenSSL](https://www.openssl.org/docs/man1.0.2/man1/openssl.html) command-line interface (CLI).

## Introduction

Welcome aboard, sailor! As Coastal Containers sails towards a secure future, you need to ensure that our ship-to-shore communications are as secure as a treasure chest. That's where SPIRE (the SPIFFE Runtime Environment) comes in. In this lab, you'll don the hat of a ship's engineer, working with pre-compiled binaries to get our SPIRE Server and Agent up and running. This is a crucial step in establishing secure, trusted communications between our fleet and the mainland.

Throughout the course of this demonstration, we'll demystify key elements of SPIRE, such as trust domains, server/agent configuration, and simple attestation methods. By the end of the hands-on exercise, not only will you have a functional SPIRE setup, but you'll also understand the 'why' and 'how' behind key configuration settings. It's not just about making things work; it's about knowing why they work.

### Preparing Your Environment

Before you cast off, prepare your ships to sail by setting your working directory in [lab-02-pki-basics](lab-02-pki-basics) as an environment variable:

```bash
export LAB_DIR=$(pwd)
```

This will make issuing commands easier in the following steps of this exercise, and will reduce the possibility of reference errors.

## Step-by-Step Instructions

### Step 1: Download SPIRE Binaries

With your working directory at [lab-02-binary-install](../lab-02-binary-install/), download the pre-built SPIRE binaries for Linux from the [SPIRE GitHub releases page](https://github.com/spiffe/spire/releases/).

```bash
export ARCH=$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)
curl -s -N -L https://github.com/spiffe/spire/releases/download/v1.8.0/spire-1.8.0-linux-$ARCH-musl.tar.gz | tar xz -C $LAB_DIR
```

This command will unpack the pre-built `spire-server` and `spire-agent` executables along with their configuration files into a directory named `spire-1.8.0` within `$LAB_DIR`.

Once unpacked, you can navigate to the directory where the SPIRE binaries were extracted.

```bash
cd $LAB_DIR/spire-1.8.0
```

Explore this directory and locate the key config files for the SPIRE Server ([server.conf](./sample/config/server/server.conf)) and Agent ([agent.conf](./sample/config/agent/agent.conf)).

### Step 2: Configure the SPIRE Server

In this step, you'll configure the SPIRE Server using an example [server.conf](./sample/config/server/server.conf) file tailored for the ControlPlane organization. We'll walk you through the key parameters you need to emulate this configuration.

To begin with server configuration, navigate to the SPIRE Server configuration directory using the `$LAB_DIR` environment variable you set previously. 

```bash
cd $LAB_DIR/spire-1.8.0/conf/server
```

Now, open your local `server.conf` file in your preferred text editor. Once opened, note the differences in the provided [server.conf](./sample/config/server/server.conf) file and update your local `server.conf` file to reflect the example config.

#### Key Parameters to Update

- `bind_address` and `bind_port`: These are set to `127.0.0.1` and `8081` for the purposes of our local demo setup. The server will listen to these when establishing a connection.
- `trust_domain`: Set this to `coastal-containers.example`. This configuration item is vital for defining the domain boundary that SPIFFE IDs can be asserted.
- Data and Key Paths: Use relative paths (e.g., `./data/..`) for `data_dir`, `keys_path`, and `connection_string` to simplify local setup in the lab repository.
- `ca_ttl`: Set to 72h (72 hours). This config sets the lifetime of the Certificate Authority (CA) signing key responsible for issuing SVIDs.
- `default_x509_svid_ttl`: Set to 6h (6 hours). This defines the X509-SVID time-to-live, or how long issued SVIDs are valid.
- `ca_subject`: This is the Subject material that CA certificates will use. For the Coastal Containers example we will be using `UK` for the `country`, `CoastalContainers` for the `organization`, and `Coastal Containers Ltd` as the `common_name`. 

An example of what this configuration will look like within the `server.conf` file is shown here:

```conf
server {
    bind_address = "127.0.0.1"
    bind_port = "8081"
    trust_domain = "coastal-containers.example"
    data_dir = "./data/server"
    log_level = "DEBUG"
    ca_ttl = "168h"
    default_x509_svid_ttl = "48h"
    
    ca_subject {
        country = ["US"]
        organization = ["CoastalContainers"]
        common_name = "Coastal Containers Ltd"
    }
}

plugins {
    DataStore "sql" {
        plugin_data {
            database_type = "sqlite3"
            connection_string = "./data/server/datastore.sqlite3"
        }
    }

    KeyManager "disk" {
        plugin_data {
            keys_path = "./data/server/keys.json"
        }
    }

    NodeAttestor "join_token" {
        plugin_data {}
    }
}
```

⚠️ *Note: These settings are for demonstration purposes only and are not suitable for production environments. For a detailed configuration guide, check the [SPIRE Server Configuration Reference](https://spiffe.io/docs/latest/deploying/spire_server/).*

### Step 3: Configure the SPIRE Agent

Now, to begin agent configuration, navigate to the SPIRE Agent configuration directory using the `$LAB_DIR` environment variable you set previously. 

```bash
cd $LAB_DIR/spire-1.8.0/conf/agent
```

You should follow a similar configuration process for the agent as you did with the server. First, open your local `agent.conf` file in your preferred text editor. Once opened, note the differences in the provided [agent.conf](./sample/config/agent/agent.conf) file and update your local `agent.conf` file to reflect the example config.

#### Key Parameters to Update

- `server_address` and `server_port`: These should match the `bind_address` and `bind_port` from the `server.conf` file, allowing the agent to locate the server.
- `trust_domain`: This should also match the `trust_domain` in the `server.conf` to maintain a consistent identity boundary.
- `insecure_bootstrap`: Set to `true` for this demo, which eases the initial agent registration with the server by allowing bootstrap without verification of the SPIRE Servers identity.
- `WorkloadAttestor`: Set to `unix`. This parameter is used for workload attestation, allowing the agent to confirm the identity of connecting workloads.
- Data and Key Paths: Similar to the server, you use relative paths (e.g., `./data/..`) for `data_dir` and the KeyManager disk `directory` to simplify the setup on your local device.

An example of what this configuration will look like within the `agent.conf` file is shown here:

```conf
agent {
    data_dir = "./data/agent"
    log_level = "DEBUG"
    trust_domain = "coastal-containers.example"
    server_address = "localhost"
    server_port = 8081

    # Insecure bootstrap is NOT appropriate for production use but is ok for 
    # simple testing/evaluation purposes.
    insecure_bootstrap = true
}

plugins {
   KeyManager "disk" {
        plugin_data {
            directory = "./data/agent"
        }
    }

    NodeAttestor "join_token" {
        plugin_data {}
    }

    WorkloadAttestor "unix" {
        plugin_data {}
    }
}

```

⚠️ *Note: These settings are for demonstration purposes only and are not suitable for production environments. For a detailed configuration guide, check the [SPIRE Agent Configuration Reference](https://spiffe.io/docs/latest/deploying/spire_agent/)*.

### Step 4: Start the SPIRE Server

Now that you have configured the server and agent, you can start running the SPIRE components. 

First, navigate back to the root SPIRE directory so you can access the `spire-server` and `spire-agent` executables.

```bash
cd $LAB_DIR/spire-1.8.0
```

Now, you are ready to start the SPIRE Server using your updated `server.conf` configuration file.

```bash
bin/spire-server run -config conf/server/server.conf &
```

#### Output Log Insights

- `INFO[0000] Configured`: indicates that the server has successfully loaded the configuration.
- `INFO[0000] X509 CA`: activated signals that the X.509 Certificate Authority has been activated, which will be responsible for issuing SVIDs.
- `INFO[0000] Starting Server APIs`: Starting Server APIs confirms that the server is starting its APIs and is ready to accept connections.
- `DEBU[0001] Initializing health checkers`: Signals that the server's health check mechanisms are initialized.

Once the output log has finished, run the following command to verify that the SPIRE Server is running and healthy.

```bash
bin/spire-server healthcheck
```

If everything is working and the SPIRE Server has started properly, you should see the following:

```bash
Server is healthy.
```

If not, you may need to troubleshoot your `server.conf` configuration and attempt the process again.

### Step 5: Generate a Join Token

For this demo, you will be using a [join token](https://spiffe.io/docs/latest/deploying/configuring/#join-token) to provide agent attestation to the SPIRE Server. A join token is a simple, one-time-use token for this attestation process. Other methods of node attestation can be found [here](https://spiffe.io/docs/latest/deploying/configuring/#how-to-configure-spire) within the official SPIFFE docs. 

Now, with your SPIRE Server up and running, generate a join token to attest the SPIRE Agent to the SPIRE Server.

```bash
bin/spire-server token generate -spiffeID spiffe://coastal-containers.example/spire-agent
```

Upon execution, you'll receive a `<token_string>`. Keep it safe—you'll need it soon.

In the context of Zero Trust, it's often said, "It's turtles all the way down," meaning that one layer of security depends on another, creating a seemingly never-ending loop. SPIRE helps you find the "bottom turtle" — the foundational secret that breaks this loop. This is particularly important in dynamic ecosystems where manual provisioning of secrets is not feasible. For the purposes of this demo, the join token serves as this foundational secret for the SPIRE agent to attest to the SPIRE server, initiating a trusted relationship without needing another secret.

*Note: The provided `spiffeID` aligns with the Coastal Container organization's URI structure. For more on SPIFFE IDs, consult [SPIFFE Concepts](https://spiffe.io/docs/latest/spiffe-about/spiffe-concepts/#spiffe-id) within the official docs.*

### Step 6: Start the SPIRE Agent

Having the join token in hand, you can now start the SPIRE Agent. The agent interacts with the SPIRE Server, receives SVIDs, and provides them to workloads. Use the join token by passing in the copied `<token_string>` to start the SPIRE Agent.

```bash
bin/spire-agent run -config conf/agent/agent.conf -joinToken <your_join_token> &
```

#### Output Log Insights

- `INFO[0000] Starting agent with data directory: "..."`: Indicates the data directory where the agent will store runtime data.
- `INFO[0000] Node attestation was successful`: Confirms that the agent successfully attested to the server. You'll also see the SPIFFE ID of the join token you generated earlier here.
- `INFO[0000] Starting Workload and SDS APIs`: Points out that the Workload API is up and ready to serve SVIDs.
- `DEBU[0001] Initializing health checkers`: Signals that the agent's health check mechanisms are initialized.

Once the output log has finished, run the following command to verify that the SPIRE Agent is running and healthy.

```bash
bin/spire-agent healthcheck
```

If everything is working and the SPIRE Agent has started properly, you should see the following:

```bash
Agent is healthy.
```

If not, you may need to troubleshoot your `agent.conf` configuration and attempt the process again.

### Step 7: Create a Registration Entry

Now that your SPIRE Server and Agent are running, you can register workloads by creating a registration entry. Registration entries are responsible for mapping a SPIFFE ID to specific selectors so workloads can be issued a valid identity. For this example, you will be creating a registration entry for the current user's `UID` as a dummy 'captain' workload.

This process requires the `unix` workload attestor configured within the `agent.conf` file.

```bash
bin/spire-server entry create -parentID spiffe://coastal-containers.example/spire-agent \
    -spiffeID spiffe://coastal-containers.example/captain-workload -selector unix:uid:$(id -u)
```

If you would like to register a different workload, you can do so by updating the `selector` field(s) and configuring the `NodeAttestor` plugin. More information on how to do this can be found in the [Registering workloads](https://spiffe.io/docs/latest/deploying/registering/) section of the official SPIFFE docs.

### Step 8: Retrieve SVID Details

To retrieve and view the details of the SVID you just created, you can use the following command to emulate the process that a workload would normally take to fetch an X.509-SVID from the SPIRE Agent.

```bash
bin/spire-agent api fetch x509 -write /tmp/
```

This will retrieve the X.509-SVID and write it to the `/tmp/` temporary directory so you can view it.

### Step 9: Inspect the SVID

Once you have written the SVID to the `/tmp/` directory, you can issue the following `openssl` command to view its contents.

```bash
openssl x509 -in /tmp/svid.0.pem -text -noout
```

Compare the created X.509-SVID certificate to the one you generated in the [initial PKI lab](../lab-01-pki-basics/). Take note of the SPIFFE ID in the URI SAN and the format of the Subject field. This is a key part of how SPIRE manages identities. In the SPIFFE ecosystem, an X.509-SVID must contain exactly one URI (SPIFFE ID) in the Subject Alternative Name (SAN) extension, as opposed to the initial X.509 Cert, which relies on geographical and organizational information in the Subject field to establish identity. 

### Step 10: Clean Up

To restore the lab directory to its original state, first change to the `$LAB_DIR` directory, and then perform the following steps:

1. Kill the SPIRE Server and Agent processes.

```bash
killall spire-server spire-agent
```

1. Delete any downloaded or generated files.

```bash
rm -rf spire-1.8.0
```

## Key Configuration Details

Understanding key configuration details is essential for anyone looking to leverage the full potential of SPIRE. In this section, we'll break down what each configuration element is designed to do, so you're not just following instructions—you're gaining a deep understanding.

### Trust Domain Configuration

The trust domain is the cornerstone of SPIRE's security model. It acts as an identity namespace and is essential for issuing and verifying SVIDs (SPIFFE Verifiable Identity Documents). As such, it should be unique to your organizational architecture and needs. In our lab, you used `coastal-containers.example` as the trust domain.

In `server.conf` and `agent.conf`:

```bash
trust_domain = "coastal-containers.example"
```

### Server Address Configuration

The `bind_address` in the `server.conf` file determines which IP address the SPIRE Server will bind to for listening to incoming connections. By default, this is set to `0.0.0.0`, meaning the server will listen on all available network interfaces. This can be changed to any specific IP address of the machine where the SPIRE Server is running.

In `server.conf`:

```bash
bind_address = "0.0.0.0"
```

In `agent.conf`:

```bash
server_address = 0.0.0.0
```

### Port Configuration

By default, the SPIRE Server listens on port 8081 for incoming connections from SPIRE Agents. If you want to change it, you can do so with the `bind_port` and `server_port` parameters.

In `server.conf`:

```bash
bind_port = "9090"
```

In `agent.conf`:

```bash
server_port = 9090
```

### Node Attestation: Using Join Token

In this lab, you're using a join token as the method of node attestation. This is a simple yet powerful way to establish initial trust between the SPIRE Server and Agent. More information about node attestation and the various node attestor plugins can be found [here](https://spiffe.io/docs/latest/spire-about/spire-concepts/#node-attestation) within the official SPIFFE docs.

In `server.conf` and `agent.conf`:

```bash
NodeAttestor "join_token" {
  plugin_data {}
}
```

Additionally, the join token will be supplied as a command-line argument when starting the agent.

### Workload Attestation: Using Unix

This lab uses the Unix Workload Attestor to verify the identity of processes. This method relies on Unix attributes like user and group IDs. More information about workload attestation and the various workload attestor plugins can be found [here](https://spiffe.io/docs/latest/spire-about/spire-concepts/#workload-attestation) within the official SPIFFE docs. 

In `agent.conf`:

```bash
WorkloadAttestor "unix" {
  plugin_data {}
}
```

### Data Storage: SQLite

The DataStore plugin provides persistent storage for the SPIRE Server. Out of the box, SPIRE supports SQLite, PostgreSQL, and MySQL databases. For the scope of this lab, you used SQLite for data storage—a lightweight, file-based database that simplifies setup.

In `server.conf`:

```bash
DataStore "sql" {
  plugin_data {
    database_type = "sqlite3"
    connection_string = "./coastal-containers/data/server/datastore.sqlite3"
  }
}
```

### Key Management: Disk

The KeyManager plugin provides signing and key storage logic for SPIRE Server signing operations. In this demo, both the SPIRE Server and Agent will store their keys locally on disk. This is a straightforward method suitable for our lab's scope, though not recommended for production environments.

In `server.conf`:

```bash
KeyManager "disk" {
  plugin_data {
    directory = "./coastal-containers/data/server/keys.json"
  }
}
```

In `agent.conf`:

```bash
KeyManager "disk" {
  plugin_data {
    directory = "./coastal-containers/data/agent"
  }
}
```

## Production Considerations

As this lab is meant for demonstration purposes only, there are a number of items that are NOT suitable for production. Due to this, we would like to address some of these configuration items at a high level to guide later labs that emulate more complex deployment scenarios for SPIRE.

### 1. `insecure_bootstrap` Option

**Consideration**: The `insecure_bootstrap` option should be set to `false` (or ideally removed) in a production environment.  

**Why**: When `insecure_bootstrap` is set to `true`, the agent skips server validation, making it vulnerable to man-in-the-middle attacks.  

**Recommended Action**: Use a secure method for the initial bootstrap of the SPIRE agent to ensure the integrity of the SPIRE server's identity.

### 2. Absolute Paths vs Relative Paths

**Consideration**: Use absolute paths instead of relative paths for directories like `./data/agent`.  

**Why**: Using absolute paths ensures that SPIRE can always locate the data directory, regardless of the current working directory from which the SPIRE services are started.  

**Recommended Action**: Update the `data_dir`, Data Store `connection_string`, and Key Manager `disk` directory configuration in both `server.conf` and `agent.conf` to use absolute paths such as `/opt/data../`

### 3. Time-to-Live (TTL) Settings

**Consideration**: Use realistic TTL settings for the CA certificate and SVIDs.  

**Why**: Long-lived certificates pose a security risk if they are compromised. Short-lived certificates reduce the window of time an attacker can use them.  

**Recommended Action**: 

- Set `default_x509_svid_ttl` and `default_jwt_svid_ttl` to a short period.
- Make sure `ca_ttl` is set to a reasonable timeframe that aligns with your CA rotation strategy.

### 4. Use Appropriate Attestors

**Consideration**: Use production-appropriate attestors for both node and workload attestation.

**Why**: The `join_token` method of attestation is easy for demos but not highly secure.

**Recommended Action**: 

- For node attestation, consider using cloud-specific attestors or TPM-based attestors.
- For workload attestation, use attestors that are suitable for your runtime environment (e.g., Kubernetes, Docker, etc.)

### 5. Secure Data Storage

**Consideration**: Use a more robust and distributed data store in production instead of SQLite.  

**Why**: SQLite is easy to set up but not designed for high-availability, failover, or distributed systems. 

**Recommended Action**: Use a production-ready database like MySQL or PostgreSQL.

### 6. Logging Levels

**Consideration**: Set an appropriate log level.  

**Why**: Debug-level logging provides detailed information but could have performance implications and may expose sensitive information in the logs.  

**Recommended Action**: Consider using log levels like `INFO` or `WARN` for production.

More information about this configuration can be found within the [Configuring SPIRE](https://spiffe.io/docs/latest/deploying/configuring/) section of the official SPIFFE docs. 

## Conclusion

Congratulations, you've just fortified Coastal Containers' communication lines, making it harder for the likes of Captain Hashjack and his band of cyber-pirates to infiltrate our systems. You've not only installed SPIRE but have also gained a deeper understanding of its inner workings. This is essential knowledge as you continue our Zero Trust voyage.

For those eager to continue sailing through uncharted waters, there are more advanced configurations and deployment strategies to explore. You can dive into more advanced [SPIRE deployment examples and configurations](https://github.com/spiffe/spire-examples). Fair winds and following seas, sailor!
