#!/bin/bash

# This script tweaks the original protodefs from carbide because they're messy in a way that keeps us from building.
# Try and keep this script safe to rerun on the protodefs multiple times.

# Cross-platform sed -i (macOS requires '', Linux doesn't)
sedi() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

# dpa_rpc.proto has a duplicate message "Metadata", we don't need any of it so just remove it
rm -f carbideproto/dpa_rpc.proto
sedi -e '/^import.*dpa_rpc/d' carbideproto/forge.proto

# nmx_c.proto and rack_manager.proto have duplicate ReturnCode/Metadata types, not needed
rm -f carbideproto/nmx_c.proto carbideproto/rack_manager.proto

# Both site explorer and main forge have a PowerState enum
sedi -e 's/ PowerState/ ComputerSystemPowerState/g' carbideproto/site_explorer.proto

# This one I'm blaming on the code generator.
sedi -e 's/MachineValidationStarted started/MachineValidationStarted oneof_started/' \
     -e 's/MachineValidationInProgress in_progress/MachineValidationInProgress oneof_in_progress/g' \
     -e 's/MachineValidationCompleted completed/MachineValidationCompleted oneof_completed/g' \
     carbideproto/forge.proto
