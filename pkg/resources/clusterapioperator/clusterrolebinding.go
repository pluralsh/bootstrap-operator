package clusterapioperator

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	v1 "k8s.io/api/rbac/v1"
)

func ManagerClusterRoleBindingCreator(data *resources.TemplateData) reconciling.NamedClusterRoleBindingCreatorGetter {
	return func() (string, reconciling.ClusterRoleBindingCreator) {
		return resources.ClusterAPIManagerClusterRoleBindingName, func(r *v1.ClusterRoleBinding) (*v1.ClusterRoleBinding, error) {
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Name = resources.ClusterAPIManagerClusterRoleBindingName
			r.RoleRef = v1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     resources.ClusterAPIManagerClusterRoleName,
			}
			r.Subjects = []v1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "default",
					Namespace: data.Namespace,
				},
			}
			return r, nil
		}
	}
}

func ProxyClusterRoleBindingCreator(data *resources.TemplateData) reconciling.NamedClusterRoleBindingCreatorGetter {
	return func() (string, reconciling.ClusterRoleBindingCreator) {
		return resources.ClusterAPIProxyClusterRoleBindingName, func(r *v1.ClusterRoleBinding) (*v1.ClusterRoleBinding, error) {
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Name = resources.ClusterAPIProxyClusterRoleBindingName
			r.RoleRef = v1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     resources.ClusterAPIProxyClusterRoleName,
			}
			r.Subjects = []v1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "default",
					Namespace: data.Namespace,
				},
			}
			return r, nil
		}
	}
}
