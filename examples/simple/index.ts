import * as k3s from "@pulumi/k3s";
import { Config } from "@pulumi/pulumi";

const cfg = new Config()

const cluster = new k3s.Cluster("mycluster", {
    masterNodes: [
        {
            host: cfg.require("master_ip"),
            privateKey: cfg.require("private_key"),
            user: "ubuntu",
        }
    ],
    versionConfig: {
        channel: "latest"
    }
});

export const kubeconfig = cluster.kubeconfig;