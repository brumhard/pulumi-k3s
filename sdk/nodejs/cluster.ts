// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import { input as inputs, output as outputs } from "./types";
import * as utilities from "./utilities";

export class Cluster extends pulumi.CustomResource {
    /**
     * Get an existing Cluster resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): Cluster {
        return new Cluster(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'k3s:index:Cluster';

    /**
     * Returns true if the given object is an instance of Cluster.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Cluster {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Cluster.__pulumiType;
    }

    public readonly agents!: pulumi.Output<outputs.Node[] | undefined>;
    public readonly cniConfig!: pulumi.Output<outputs.CNIConfiguration | undefined>;
    public /*out*/ readonly kubeconfig!: pulumi.Output<string>;
    public readonly masterNodes!: pulumi.Output<outputs.Node[]>;
    public readonly versionConfig!: pulumi.Output<outputs.VersionConfiguration | undefined>;

    /**
     * Create a Cluster resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: ClusterArgs, opts?: pulumi.CustomResourceOptions) {
        let inputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.masterNodes === undefined) && !opts.urn) {
                throw new Error("Missing required property 'masterNodes'");
            }
            inputs["agents"] = args ? args.agents : undefined;
            inputs["cniConfig"] = args ? args.cniConfig : undefined;
            inputs["masterNodes"] = args ? args.masterNodes : undefined;
            inputs["versionConfig"] = args ? args.versionConfig : undefined;
            inputs["kubeconfig"] = undefined /*out*/;
        } else {
            inputs["agents"] = undefined /*out*/;
            inputs["cniConfig"] = undefined /*out*/;
            inputs["kubeconfig"] = undefined /*out*/;
            inputs["masterNodes"] = undefined /*out*/;
            inputs["versionConfig"] = undefined /*out*/;
        }
        if (!opts.version) {
            opts = pulumi.mergeOptions(opts, { version: utilities.getVersion()});
        }
        super(Cluster.__pulumiType, name, inputs, opts);
    }
}

/**
 * The set of arguments for constructing a Cluster resource.
 */
export interface ClusterArgs {
    agents?: pulumi.Input<pulumi.Input<inputs.NodeArgs>[]>;
    cniConfig?: pulumi.Input<inputs.CNIConfigurationArgs>;
    masterNodes: pulumi.Input<pulumi.Input<inputs.NodeArgs>[]>;
    versionConfig?: pulumi.Input<inputs.VersionConfigurationArgs>;
}
