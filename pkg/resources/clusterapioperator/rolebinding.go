package clusterapioperator

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	v1 "k8s.io/api/rbac/v1"
)

func LeaderElectionRoleBindingCreator(data *resources.TemplateData) reconciling.NamedRoleBindingCreatorGetter {
	return func() (string, reconciling.RoleBindingCreator) {
		return resources.ClusterAPILeaderElectionRoleBindingName, func(r *v1.RoleBinding) (*v1.RoleBinding, error) {
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Name = resources.ClusterAPILeaderElectionRoleBindingName
			r.Namespace = data.Namespace
			r.RoleRef = v1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     resources.ClusterAPILeaderElectionRoleName,
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
