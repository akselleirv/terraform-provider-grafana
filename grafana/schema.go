package grafana

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// schemaDiffFloat32 is a SchemaDiffSuppressFunc for diffing float32 values.
// schema.TypeFloat uses float64, which is a problem for API types that use
// float32. Terraform automatically converts float32 to float64 which changes
// the precision and causes incorrect diffs.
//
// For example, synthetic_monitoring.Probe.Latitude is float32. Attempting to
// set grafanacloud_synthetic_monitoring_probe.latitude to 27.98606 results in
// 27.986059188842773. The solution is to diff old and new values as float32.
func schemaDiffFloat32(k, old string, nw string, d *schema.ResourceData) bool {
	old32, _ := strconv.ParseFloat(old, 32)
	nw32, _ := strconv.ParseFloat(nw, 32)
	return old32 == nw32
}

func cloneResourceSchemaForDatasource(r *schema.Resource, updates map[string]*schema.Schema) map[string]*schema.Schema {
	resourceSchema := r.Schema
	clone := make(map[string]*schema.Schema)
	for k, v := range resourceSchema {
		clone[k] = v
		clone[k].Computed = true
		clone[k].Optional = false
		clone[k].Required = false
		clone[k].Default = nil
		clone[k].StateFunc = nil
		clone[k].DiffSuppressFunc = nil
		clone[k].ValidateDiagFunc = nil
		clone[k].ValidateFunc = nil
	}
	for k, v := range updates {
		if v == nil {
			delete(clone, k)
		} else {
			clone[k] = v
		}
	}
	return clone
}

func allowedValuesDescription(description string, allowedValues []string) string {
	return fmt.Sprintf("%s. Allowed values: `%s`.", description, strings.Join(allowedValues, "`, `"))
}
