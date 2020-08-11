/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore 
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package agent
// ListBackups struct for ListBackups
type ListBackups struct {
	Size int32 `json:"size,omitempty"`
	Items []Backup `json:"items,omitempty"`
}
