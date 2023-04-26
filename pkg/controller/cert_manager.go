package controller

import (
	"context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Reconciler) CheckCertManager(ctx context.Context) bool {
	cert := &cmapi.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cmapichecker",
			Namespace: r.Namespace,
		},
		Spec: cmapi.CertificateSpec{
			DNSNames:   []string{"cmapichecker.example"},
			SecretName: "cmapichecker",
			IssuerRef: cmmeta.ObjectReference{
				Name: "cmapichecker",
			},
		},
	}
	if err := r.Create(ctx, cert); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return true
		}
		return false
	}
	return true
}
