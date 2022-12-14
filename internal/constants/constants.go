// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package constants defines some global constants
package constants

const (
	// ControllerContainerName is the name of Controller container
	ControllerContainerName = "controller"

	// ControllerConfigMapName is the name of Controller configmap
	ControllerConfigMapName = "config-controller"

	// StoreContainerName is the name of store container
	StoreContainerName = "store"

	// BrokerConfigName is the name of mounted configuration file
	StoreConfigMapName = "config-store"

	// TriggerContainerName is the name of tigger container
	TriggerContainerName = "tigger"

	// BrokerConfigName is the name of mounted configuration file
	TriggerConfigMapName = "config-trigger"

	// TimerContainerName is the name of timer container
	TimerContainerName = "timer"

	// BrokerConfigName is the name of mounted configuration file
	TimerConfigMapName = "config-timer"

	// GatewayContainerName is the name of gateway container
	GatewayContainerName = "gateway"

	// BrokerConfigName is the name of mounted configuration file
	GatewayConfigMapName = "config-gateway"

	// ConfigMountPath is the directory of Vanus configd files
	ConfigMountPath = "/vanus/config"

	// VolumeMountPath is the directory of Store data files
	VolumeMountPath = "/data"

	// VolumeName is the directory of Store data files
	VolumeName = "data"

	// VolumeName is the directory of Store data files
	VolumeStorage = "1Gi"

	// StorageModeStorageClass is the name of StorageClass storage mode
	StorageModeStorageClass = "StorageClass"

	// StorageModeEmptyDir is the name of EmptyDir storage mode
	StorageModeEmptyDir = "EmptyDir"

	// StorageModeHostPath is the name pf HostPath storage mode
	StorageModeHostPath = "HostPath"

	// the container environment variable name of controller pod ip
	EnvPodIP = "POD_IP"

	// the container environment variable name of controller pod name
	EnvPodName = "POD_NAME"

	// the container environment variable name of controller log level
	EnvLogLevel = "VANUS_LOG_LEVEL"
)

const (
	// ServiceAccountName is the ServiceAccount name of Vanus cluster
	ServiceAccountName = "vanus-operator"

	ContainerPortNameGrpc        = "grpc"
	ContainerPortNameEtcdClient  = "etcd-client"
	ContainerPortNameEtcdPeer    = "etcd-peer"
	ContainerPortNameMetrics     = "metrics"
	ContainerPortNameProxy       = "proxy"
	ContainerPortNameCloudevents = "cloudevents"

	ControllerPortGrpc       = 2048
	ControllerPortEtcdClient = 2379
	ControllerPortEtcdPeer   = 2380
	ControllerPortMetrics    = 2112

	StorePortGrpc          = 11811
	TriggerPortGrpc        = 2148
	GatewayPortProxy       = 8080
	GatewayPortCloudevents = 8081

	HeadlessService = "None"

	// RequeueIntervalInSecond is an universal interval of the reconcile function
	RequeueIntervalInSecond = 6
)

var (
	// GroupNum is the number of broker group
	GroupNum = 0

	// NameServersStr is the name server list
	NameServersStr = ""

	// IsNameServersStrUpdated is whether the name server list is updated
	IsNameServersStrUpdated = false

	// IsNameServersStrInitialized is whether the name server list is initialized
	IsNameServersStrInitialized = false

	// BrokerClusterName is the broker cluster name
	BrokerClusterName = ""

	// svc of controller for brokers
	ControllerAccessPoint = ""
)
