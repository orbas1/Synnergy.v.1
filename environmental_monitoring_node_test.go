package synnergy

import "testing"

func TestEnvironmentalMonitoringNodeTrigger(t *testing.T) {
	n := NewEnvironmentalMonitoringNode()
	cond := EnvCondition{SensorID: "temp1", Operator: ">", Threshold: 50}
	n.SetCondition(cond)
	if !n.Trigger("temp1", []byte("55")) {
		t.Fatal("expected trigger for value above threshold")
	}
	if n.Trigger("temp1", []byte("45")) {
		t.Fatal("did not expect trigger for value below threshold")
	}
}
