# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities
from . import outputs

__all__ = [
    'CRIConfiguration',
    'Node',
    'VersionConfiguration',
]

@pulumi.output_type
class CRIConfiguration(dict):
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "enableGVisor":
            suggest = "enable_g_visor"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in CRIConfiguration. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        CRIConfiguration.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        CRIConfiguration.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 enable_g_visor: Optional[bool] = None):
        if enable_g_visor is not None:
            pulumi.set(__self__, "enable_g_visor", enable_g_visor)

    @property
    @pulumi.getter(name="enableGVisor")
    def enable_g_visor(self) -> Optional[bool]:
        return pulumi.get(self, "enable_g_visor")


@pulumi.output_type
class Node(dict):
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "privateKey":
            suggest = "private_key"
        elif key == "criConfig":
            suggest = "cri_config"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in Node. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        Node.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        Node.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 host: str,
                 private_key: str,
                 args: Optional[Sequence[str]] = None,
                 cri_config: Optional['outputs.CRIConfiguration'] = None,
                 user: Optional[str] = None):
        pulumi.set(__self__, "host", host)
        pulumi.set(__self__, "private_key", private_key)
        if args is not None:
            pulumi.set(__self__, "args", args)
        if cri_config is not None:
            pulumi.set(__self__, "cri_config", cri_config)
        if user is None:
            user = 'root'
        if user is not None:
            pulumi.set(__self__, "user", user)

    @property
    @pulumi.getter
    def host(self) -> str:
        return pulumi.get(self, "host")

    @property
    @pulumi.getter(name="privateKey")
    def private_key(self) -> str:
        return pulumi.get(self, "private_key")

    @property
    @pulumi.getter
    def args(self) -> Optional[Sequence[str]]:
        return pulumi.get(self, "args")

    @property
    @pulumi.getter(name="criConfig")
    def cri_config(self) -> Optional['outputs.CRIConfiguration']:
        return pulumi.get(self, "cri_config")

    @property
    @pulumi.getter
    def user(self) -> Optional[str]:
        return pulumi.get(self, "user")


@pulumi.output_type
class VersionConfiguration(dict):
    def __init__(__self__, *,
                 channel: Optional[str] = None,
                 version: Optional[str] = None):
        if channel is not None:
            pulumi.set(__self__, "channel", channel)
        if version is not None:
            pulumi.set(__self__, "version", version)

    @property
    @pulumi.getter
    def channel(self) -> Optional[str]:
        return pulumi.get(self, "channel")

    @property
    @pulumi.getter
    def version(self) -> Optional[str]:
        return pulumi.get(self, "version")


