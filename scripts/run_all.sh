#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SELF="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/$(basename "${BASH_SOURCE[0]}")"

for script in "$SCRIPT_DIR"/*.sh; do
    [[ "$script" == "$SELF" ]] && continue

    if [[ -f "$script" && -x "$script" ]]; then
        echo "Running: $(basename "$script")"
        "$script"
    fi
done