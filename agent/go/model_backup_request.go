/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type BackupRequest struct {
	Bucket string `json:"bucket"`

	Location string `json:"location,omitempty"`

	Envs []EnvVar `json:"envs,omitempty"`
}
