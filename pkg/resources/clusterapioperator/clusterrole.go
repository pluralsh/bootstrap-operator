package clusterapioperator

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	v1 "k8s.io/api/rbac/v1"
)

func ManagerClusterRoleCreator() reconciling.NamedClusterRoleCreatorGetter {
	return func() (string, reconciling.ClusterRoleCreator) {
		return resources.ClusterAPIManagerClusterRoleName, func(r *v1.ClusterRole) (*v1.ClusterRole, error) {
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Name = resources.ClusterAPIManagerClusterRoleName
			r.Rules = []v1.PolicyRule{
				{
					Verbs:     []string{"*"},
					APIGroups: []string{"*"},
					Resources: []string{"*"},
				},
			}

			return r, nil
		}
	}
}

func MetricsClusterRoleCreator() reconciling.NamedClusterRoleCreatorGetter {
	return func() (string, reconciling.ClusterRoleCreator) {
		return resources.ClusterAPIMetricsClusterRoleName, func(r *v1.ClusterRole) (*v1.ClusterRole, error) {
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Name = resources.ClusterAPIMetricsClusterRoleName
			r.Rules = []v1.PolicyRule{
				{
					Verbs:           []string{"get"},
					NonResourceURLs: []string{"/metrics"},
				},
			}

			return r, nil
		}
	}
}

func ProxyClusterRoleCreator() reconciling.NamedClusterRoleCreatorGetter {
	return func() (string, reconciling.ClusterRoleCreator) {
		return resources.ClusterAPIProxyClusterRoleName, func(r *v1.ClusterRole) (*v1.ClusterRole, error) {
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Name = resources.ClusterAPIProxyClusterRoleName
			r.Rules = []v1.PolicyRule{
				{
					Verbs:     []string{"create"},
					APIGroups: []string{"authentication.k8s.io"},
					Resources: []string{"tokenreviews"},
				},
				{
					Verbs:     []string{"create"},
					APIGroups: []string{"authorization.k8s.io"},
					Resources: []string{"subjectaccessreviews"},
				},
			}

			return r, nil
		}
	}
}
