package controller

import (
	"github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	extfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/kubernetes/pkg/controller"
	"github.com/longhorn/longhorn-manager/util"
	lhfake "github.com/longhorn/longhorn-manager/k8s/pkg/client/clientset/versioned/fake"
)

func (s *TestSuite) TestKubernetesPodControllerStructuredLogging(c *C) {
	lhClient := lhfake.NewSimpleClientset()
	kubeClient := k8sfake.NewSimpleClientset()
	extensionsClient := extfake.NewSimpleClientset()
	informerFactories := util.NewInformerFactories(TestNamespace, kubeClient, lhClient, controller.NoResyncPeriodFunc())
	
	kc, err := newTestKubernetesPodController(lhClient, kubeClient, extensionsClient, informerFactories)
	c.Assert(err, IsNil)
	c.Assert(kc.logger, NotNil)
	
	// Verify that kc.logger is a FieldLogger (supports structured logging)
	// The controller uses a *logrus.Entry internally after wrapping by newBaseControllerWithQueue
	c.Assert(kc.logger, FitsTypeOf, (*logrus.Entry)(nil))
	
	// Verify that kc.logger.WithError() method exists (structured logging)
	loggerWithErr := kc.logger.WithError(nil)
	c.Assert(loggerWithErr, NotNil)
	
	// Verify that kc.logger.WithField() method exists (structured logging)
	loggerWithField := kc.logger.WithField("test", "value")
	c.Assert(loggerWithField, NotNil)
	
	// Verify that the logger can log messages
	loggerWithField.Debugf("test message")
	loggerWithField.Infof("test message")
	loggerWithField.Warnf("test message")
	loggerWithField.Errorf("test message")
}
