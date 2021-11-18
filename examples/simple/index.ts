import * as k3s from "@pulumi/k3s";
import { Config, secret } from "@pulumi/pulumi";

const cfg = new Config()

const cluster = new k3s.Cluster("mycluster", {
    masterNodes: [
        {
            host: cfg.require("master_ip"),
            privateKey: cfg.requireSecret("private_key"),
            user: "ubuntu",
        }
    ],
    versionConfig: {
        channel: "latest"
    }
});

export const kubeconfig = secret(cluster.kubeconfig);