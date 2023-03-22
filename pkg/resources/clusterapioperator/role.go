package clusterapioperator

import (
	"github.com/pluralsh/bootstrap-operator/pkg/resources"
	"github.com/pluralsh/bootstrap-operator/pkg/resources/reconciling"
	v1 "k8s.io/api/rbac/v1"
)

func LeaderElectionRoleCreator(data *resources.TemplateData) reconciling.NamedRoleCreatorGetter {
	return func() (string, reconciling.RoleCreator) {
		return resources.ClusterAPILeaderElectionRoleName, func(r *v1.Role) (*v1.Role, error) {
			r.Name = resources.ClusterAPILeaderElectionRoleName
			r.Namespace = data.Namespace
			r.Labels = map[string]string{"clusterctl.cluster.x-k8s.io/core": "capi-operator"}
			r.Rules = []v1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
					APIGroups: []string{""},
					Resources: []string{"configmaps"},
				},
				{
					Verbs:     []string{"get", "update", "patch"},
					APIGroups: []string{""},
					Resources: []string{"configmaps/status"},
				},
				{
					Verbs:     []string{"create"},
					APIGroups: []string{""},
					Resources: []string{"events"},
				},
				{
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
					APIGroups: []string{"coordination.k8s.io"},
					Resources: []string{"leases"},
				},
			}

			return r, nil
		}
	}
}
