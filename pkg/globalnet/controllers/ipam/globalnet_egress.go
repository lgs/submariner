package ipam

import (
	"fmt"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/submariner-io/admiral/pkg/log"
	"k8s.io/klog"

	"github.com/submariner-io/submariner/pkg/routeagent/constants"
)

func (i *Controller) updateEgressRulesForResource(resourceName, sourceIP, globalIP string, addRules bool) error {
	ruleSpec := []string{"-p", "all", "-s", sourceIP, "-m", "mark", "--mark", globalNetIPTableMark, "-j", "SNAT", "--to", globalIP}

	if addRules {
		klog.V(log.DEBUG).Infof("Installing iptable egress rules for %s: %s", resourceName, strings.Join(ruleSpec, " "))

		if err := i.ipt.AppendUnique("nat", constants.SmGlobalnetEgressChain, ruleSpec...); err != nil {
			return fmt.Errorf("error appending iptables rule \"%s\": %v\n", strings.Join(ruleSpec, " "), err)
		}
	} else {
		klog.V(log.DEBUG).Infof("Deleting iptable egress rules for %s : %s", resourceName, strings.Join(ruleSpec, " "))
		if err := i.ipt.Delete("nat", constants.SmGlobalnetEgressChain, ruleSpec...); err != nil {
			return fmt.Errorf("error deleting iptables rule \"%s\": %v\n", strings.Join(ruleSpec, " "), err)
		}
	}

	return nil
}

func MarkRemoteClusterTraffic(ipt *iptables.IPTables, remoteCidr string, addRules bool) {
	ruleSpec := []string{"-d", remoteCidr, "-j", "MARK", "--set-mark", globalNetIPTableMark}

	if addRules {
		klog.V(log.DEBUG).Infof("Marking traffic destined to remote cluster: %s", strings.Join(ruleSpec, " "))

		if err := ipt.AppendUnique("nat", constants.SmGlobalnetMarkChain, ruleSpec...); err != nil {
			klog.Errorf("error appending iptables rule \"%s\": %v\n", strings.Join(ruleSpec, " "), err)
		}
	} else {
		klog.V(log.DEBUG).Infof("Deleting rule that marks remote cluster traffic: %s", strings.Join(ruleSpec, " "))
		if err := ipt.Delete("nat", constants.SmGlobalnetMarkChain, ruleSpec...); err != nil {
			klog.Errorf("error deleting iptables rule \"%s\": %v\n", strings.Join(ruleSpec, " "), err)
		}
	}
}
