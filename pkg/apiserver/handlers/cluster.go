// Copyright 2017 EasyStack, Inc.
package handlers

import (
	stderr "errors"
	"fmt"

	"github.com/linkall-labs/vanus-operator/api/models"
	"github.com/linkall-labs/vanus-operator/api/restapi/operations/cluster"
	cons "github.com/linkall-labs/vanus-operator/internal/constants"
	"github.com/linkall-labs/vanus-operator/pkg/apiserver/utils"

	"github.com/go-openapi/runtime/middleware"
	vanusv1alpha1 "github.com/linkall-labs/vanus-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	log "k8s.io/klog/v2"
)

var (
	WorkloadOldReplicAnno = "replic.ecns.io/workload"
	// revisionAnnoKey       = "deployment.kubernetes.io/revision"
	// rolloutReasonAnno     = "kubernetes.io/change-cause"
)

// 所有注册的处理函数 应该 按照顺序出现在 Registxxx 下

func RegistClusterHandler(a *Api) {
	a.ClusterCreateClusterHandler = cluster.CreateClusterHandlerFunc(a.createClusterHandler)
	a.ClusterDeleteClusterHandler = cluster.DeleteClusterHandlerFunc(a.deleteClusterHandler)
	a.ClusterPatchClusterHandler = cluster.PatchClusterHandlerFunc(a.patchClusterHandler)
	a.ClusterGetClusterHandler = cluster.GetClusterHandlerFunc(a.getClusterHandler)
}

/*
创建 cluster
1 参数检查，version：如果有，则部署指定镜像版本，如果没有，默认最新版本，storage size：如果没有，默认1G
2 创建各组件资源, 创建各组件crd
*/
func (a *Api) createClusterHandler(params cluster.CreateClusterParams) middleware.Responder {
	// Parse cluster params
	c, err := genClusterConfig(params)
	if err != nil {
		log.Error(err, "parse cluster params failed")
		return utils.Response(0, err)
	}

	log.Infof("parse cluster params finish, version: %s, namespace: %s, controller_replicas: %d, controller_storage_size: %s, store_replicas: %d, store_storage_size: %s\n", c.version, c.namespace, c.controllerReplicas, c.controllerStorageSize, c.storeReplicas, c.storeStorageSize)

	// Check if the cluster already exists, if exist, return error
	exist, err := a.checkClusterExist(c)
	if err != nil {
		log.Error(err, "check cluster exist failed")
		return utils.Response(0, err)
	}
	if exist {
		log.Warning("Cluster already exist")
		return utils.Response(0, stderr.New("cluster already exist"))
	}

	log.Info("Creating a new Controller.", "Controller.Namespace", c.namespace, "Controller.Name", cons.DefaultControllerName)
	controller := generateController(c)
	err = a.ctrl.Create(controller)
	if err != nil {
		log.Error(err, "Failed to create new Controller", "Controller.Namespace", c.namespace, "Controller.Name", cons.DefaultControllerName)
		return utils.Response(0, err)
	} else {
		log.Info("Successfully create Controller")
	}

	log.Info("Creating a new Store.", "Store.Namespace", c.namespace, "Store.Name", cons.DefaultStoreName)
	store := generateStore(c)
	err = a.ctrl.Create(store)
	if err != nil {
		log.Error(err, "Failed to create new Store", "Store.Namespace", c.namespace, "Store.Name", cons.DefaultStoreName)
		return utils.Response(0, err)
	} else {
		log.Info("Successfully create Store")
	}

	log.Info("Creating a new Trigger.", "Trigger.Namespace", c.namespace, "Trigger.Name", cons.DefaultTriggerName)
	trigger := generateTrigger(c)
	err = a.ctrl.Create(trigger)
	if err != nil {
		log.Error(err, "Failed to create new Trigger", "Trigger.Namespace", c.namespace, "Trigger.Name", cons.DefaultTriggerName)
		return utils.Response(0, err)
	} else {
		log.Info("Successfully create Trigger")
	}

	log.Info("Creating a new Timer.", "Timer.Namespace", c.namespace, "Timer.Name", cons.DefaultTimerName)
	timer := generateTimer(c)
	err = a.ctrl.Create(timer)
	if err != nil {
		log.Error(err, "Failed to create new Timer", "Timer.Namespace", c.namespace, "Timer.Name", cons.DefaultTimerName)
		return utils.Response(0, err)
	} else {
		log.Info("Successfully create Timer")
	}

	log.Info("Creating a new Gateway.", "Gateway.Namespace", c.namespace, "Gateway.Name", cons.DefaultGatewayName)
	gateway := generateGateway(c)
	err = a.ctrl.Create(gateway)
	if err != nil {
		log.Error(err, "Failed to create new Gateway", "Gateway.Namespace", c.namespace, "Gateway.Name", cons.DefaultGatewayName)
		return utils.Response(0, err)
	} else {
		log.Info("Successfully create Gateway")
	}

	return cluster.NewCreateClusterOK().WithPayload(nil)
}

/*
删除 cluster
1 若有随删除而删除pvc，则删除pvc
2 删除deployemnt
3 删除 vm cr
*/

func (a *Api) deleteClusterHandler(params cluster.DeleteClusterParams) middleware.Responder {
	return cluster.NewDeleteClusterOK().WithPayload(nil)
}

/*
更新 cluster
- label修改，替换+添加+删除
- 手动伸缩
- 扩容策略
*/
func (a *Api) patchClusterHandler(params cluster.PatchClusterParams) middleware.Responder {
	return cluster.NewPatchClusterOK().WithPayload(nil)
}

/*
查询 cluster
- 集群拓扑
- 各组件运行状态
- 规格
- 等等
*/
func (a *Api) getClusterHandler(params cluster.GetClusterParams) middleware.Responder {

	// containers := make([]corev1.Container, 1)
	// containers[0] = corev1.Container{
	// 	Name:  "container1",
	// 	Image: "nginx:latest",
	// }
	// creatPod := &corev1.Pod{
	// 	Spec: corev1.PodSpec{
	// 		Containers: containers,
	// 	},
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      "nginx",
	// 		Namespace: "default",
	// 	},
	// }
	// err := a.ctrl.Create(creatPod)
	// if err != nil {
	// 	klog.Infof("create pod failed:%v", err)
	// 	return cluster.NewGetClusterOK().WithPayload(nil)
	// }

	retcode := int32(400)
	msg := "niubiplus"
	return cluster.NewGetClusterOK().WithPayload(&cluster.GetClusterOKBody{
		Code: &retcode,
		Data: &models.ClusterInfo{
			Version: "v0.5.7",
		},
		Message: &msg,
	})
	// return cluster.NewGetClusterOK().WithPayload(nil)
}

func (a *Api) checkClusterExist(c *config) (bool, error) {
	err := a.ctrl.Get(types.NamespacedName{Name: cons.DefaultControllerName, Namespace: c.namespace}, &vanusv1alpha1.Controller{})
	if err == nil {
		log.Warning("Cluster already exist")
		return true, nil
	}
	if errors.IsNotFound(err) {
		return false, nil
	} else {
		log.Error(err, "Failed to get Controller.")
		return false, err
	}
	// TODO(jiangkai): need to check other components
}

type config struct {
	namespace             string
	version               string
	controllerReplicas    int32
	controllerStorageSize string
	storeReplicas         int32
	storeStorageSize      string
}

func genClusterConfig(params cluster.CreateClusterParams) (*config, error) {
	// check required parameters
	if params.Cluster.Version == "" {
		return nil, stderr.New("cluster version is required parameters")
	}
	c := &config{
		version:               params.Cluster.Version,
		namespace:             cons.DefaultNamespace,
		controllerReplicas:    cons.DefaultControllerReplicas,
		controllerStorageSize: cons.VolumeStorage,
		storeReplicas:         cons.DefaultStoreReplicas,
		storeStorageSize:      cons.VolumeStorage,
	}
	if params.Cluster.Namespace != "" {
		c.namespace = params.Cluster.Namespace
	}
	if params.Cluster.ControllerReplicas != 0 && params.Cluster.ControllerReplicas != cons.DefaultControllerReplicas {
		c.controllerReplicas = params.Cluster.ControllerReplicas
	}
	if params.Cluster.ControllerStorageSize != "" {
		c.controllerStorageSize = params.Cluster.ControllerStorageSize
	}
	if params.Cluster.StoreReplicas != 0 && params.Cluster.StoreReplicas != cons.DefaultStoreReplicas {
		c.storeReplicas = params.Cluster.StoreReplicas
	}
	if params.Cluster.ControllerStorageSize != "" {
		c.storeStorageSize = params.Cluster.StoreStorageSize
	}
	return c, nil
}

func labelsForController(name string) map[string]string {
	return map[string]string{"app": name}
}

func generateController(c *config) *vanusv1alpha1.Controller {
	labels := labelsForController(cons.DefaultControllerName)
	requests := make(map[corev1.ResourceName]resource.Quantity)
	requests[corev1.ResourceStorage] = resource.MustParse(c.controllerStorageSize)
	controller := &vanusv1alpha1.Controller{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      cons.DefaultControllerName,
		},
		Spec: vanusv1alpha1.ControllerSpec{
			Replicas:        &c.controllerReplicas,
			Image:           fmt.Sprintf("public.ecr.aws/vanus/controller:%s", c.version),
			ImagePullPolicy: corev1.PullIfNotPresent,
			Resources:       corev1.ResourceRequirements{},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Name:   cons.VolumeName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: requests,
					},
				},
			}},
		},
	}
	return controller
}

func generateStore(c *config) *vanusv1alpha1.Store {
	labels := labelsForController(cons.DefaultStoreName)
	requests := make(map[corev1.ResourceName]resource.Quantity)
	requests[corev1.ResourceStorage] = resource.MustParse(c.storeStorageSize)
	store := &vanusv1alpha1.Store{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      cons.DefaultStoreName,
		},
		Spec: vanusv1alpha1.StoreSpec{
			Replicas:        &c.storeReplicas,
			Image:           fmt.Sprintf("public.ecr.aws/vanus/store:%s", c.version),
			ImagePullPolicy: corev1.PullIfNotPresent,
			Resources:       corev1.ResourceRequirements{},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Name:   cons.VolumeName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: requests,
					},
				},
			}},
		},
	}
	return store
}

func generateTrigger(c *config) *vanusv1alpha1.Trigger {
	trigger := &vanusv1alpha1.Trigger{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      cons.DefaultTriggerName,
		},
		Spec: vanusv1alpha1.TriggerSpec{
			Image:           fmt.Sprintf("public.ecr.aws/vanus/trigger:%s", c.version),
			ImagePullPolicy: corev1.PullIfNotPresent,
		},
	}
	return trigger
}

func generateTimer(c *config) *vanusv1alpha1.Timer {
	timer := &vanusv1alpha1.Timer{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      cons.DefaultTimerName,
		},
		Spec: vanusv1alpha1.TimerSpec{
			Image:           fmt.Sprintf("public.ecr.aws/vanus/timer:%s", c.version),
			ImagePullPolicy: corev1.PullIfNotPresent,
		},
	}
	return timer
}

func generateGateway(c *config) *vanusv1alpha1.Gateway {
	gateway := &vanusv1alpha1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      cons.DefaultGatewayName,
		},
		Spec: vanusv1alpha1.GatewaySpec{
			Image:           fmt.Sprintf("public.ecr.aws/vanus/gateway:%s", c.version),
			ImagePullPolicy: corev1.PullIfNotPresent,
		},
	}
	return gateway
}
