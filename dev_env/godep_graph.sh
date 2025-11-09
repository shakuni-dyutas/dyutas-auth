#!/bin/bash

# Generate mermaid dependency graph
output=$(
    godepgraph \
        -o ./cmd,./internal/rest,github.com/shakuni-dyutas/dyutas-auth \
        -i github.com/shakuni-dyutas/dyutas-auth/internal/helper \
        ./internal/rest
)

[ -z "$output" ] && echo "Error: godepgraph produced no output." >&2 && exit 1

output=$(
    echo "$output" | sed 's|\github.com/shakuni-dyutas/dyutas-auth/internal/||g' \
    | sed 's|\./cmd/||g' \
    | sed 's|\./internal/||g'
)

dot -Tpng -o pkg_dep_graph.png <(echo "$output")

echo "Package dependency graph generated and saved to pkg_dep_graph.png"