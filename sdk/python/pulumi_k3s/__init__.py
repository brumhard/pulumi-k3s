# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

from . import _utilities
import typing
# Export this package's modules as members:
from .cluster import *
from .provider import *
_utilities.register(
    resource_modules="""
[
 {
  "pkg": "k3s",
  "mod": "index",
  "fqn": "pulumi_k3s",
  "classes": {
   "k3s:index:Cluster": "Cluster"
  }
 }
]
""",
    resource_packages="""
[
 {
  "pkg": "k3s",
  "token": "pulumi:providers:k3s",
  "fqn": "pulumi_k3s",
  "class": "Provider"
 }
]
"""
)
