// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.K3s.Inputs
{

    public sealed class CRIConfigurationArgs : global::Pulumi.ResourceArgs
    {
        [Input("enableGVisor")]
        public Input<bool>? EnableGVisor { get; set; }

        public CRIConfigurationArgs()
        {
        }
        public static new CRIConfigurationArgs Empty => new CRIConfigurationArgs();
    }
}
