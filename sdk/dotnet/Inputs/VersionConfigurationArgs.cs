// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.K3s.Inputs
{

    public sealed class VersionConfigurationArgs : global::Pulumi.ResourceArgs
    {
        [Input("channel")]
        public Input<string>? Channel { get; set; }

        [Input("version")]
        public Input<string>? Version { get; set; }

        public VersionConfigurationArgs()
        {
        }
        public static new VersionConfigurationArgs Empty => new VersionConfigurationArgs();
    }
}
