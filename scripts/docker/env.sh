# App info
export app_name=toes
export version=$(git describe --tags --match='v*' | sed 's/^v//')

# Build
export registry_prefix=toes
export images=(toes-apiserver)
export architecture=amd64
