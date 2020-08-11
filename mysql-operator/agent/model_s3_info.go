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

// S3Info struct for S3Info
type S3Info struct {
	Bucket    string    `json:"bucket"`
	Path      string    `json:"path,omitempty"`
	AwsConfig AwsConfig `json:"awsConfig,omitempty"`
}
