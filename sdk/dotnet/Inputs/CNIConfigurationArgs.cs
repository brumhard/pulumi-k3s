// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.K3s.Inputs
{

    public sealed class CNIConfigurationArgs : Pulumi.ResourceArgs
    {
        [Input("provider")]
        public Input<string>? Provider { get; set; }

        public CNIConfigurationArgs()
        {
        }
    }
}