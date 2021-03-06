package framework

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/catalog/v1alpha1"
	. "github.com/onsi/gomega"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *Invocation) MySQLVersion() *api.MySQLVersion {
	return &api.MySQLVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: DBCatalogName,
			Labels: map[string]string{
				"app": i.app,
			},
		},
		Spec: api.MySQLVersionSpec{
			Version: DBVersion,
			DB: api.MySQLVersionDatabase{
				Image: fmt.Sprintf("%s/mysql:%s", DockerRegistry, DBVersion),
			},
			Exporter: api.MySQLVersionExporter{
				Image: fmt.Sprintf("%s/mysqld-exporter:%s", DockerRegistry, ExporterTag),
			},
			Tools: api.MySQLVersionTools{
				Image: fmt.Sprintf("%s/mysql-tools:%s", DockerRegistry, DBToolsTag),
			},
		},
	}
}

func (f *Framework) CreateMySQLVersion(obj *api.MySQLVersion) error {
	_, err := f.extClient.CatalogV1alpha1().MySQLVersions().Create(obj)
	if err != nil && kerr.IsAlreadyExists(err) {
		e2 := f.extClient.CatalogV1alpha1().MySQLVersions().Delete(obj.Name, &metav1.DeleteOptions{})
		Expect(e2).NotTo(HaveOccurred())
		_, e2 = f.extClient.CatalogV1alpha1().MySQLVersions().Create(obj)
		return e2
	}
	return nil
}

func (f *Framework) DeleteMySQLVersion(meta metav1.ObjectMeta) error {
	return f.extClient.CatalogV1alpha1().MySQLVersions().Delete(meta.Name, &metav1.DeleteOptions{})
}
