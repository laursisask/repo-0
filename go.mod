module github.com/keymone/kind

go 1.13

require (
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.3
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	gopkg.in/yaml.v3 v3.0.0-20191010095647-fc94e3f71652
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	sigs.k8s.io/kind v0.5.1
	sigs.k8s.io/yaml v1.1.0
)

replace sigs.k8s.io/kind => ./
