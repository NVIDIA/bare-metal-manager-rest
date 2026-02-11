package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nvidia/carbide-rest/rla/pkg/client"
	"github.com/nvidia/carbide-rest/rla/pkg/types"
)

var (
	actualCmd = &cobra.Command{
		Use:   "actual",
		Short: "Get actual components from external systems (e.g., Carbide)",
		Long: `Get actual components from external systems like Carbide.

Currently only supports Compute component type.

Specify exactly ONE of the following options:
  --rack-ids      : Comma-separated list of rack UUIDs
  --rack-names    : Comma-separated list of rack names
  --component-ids : Comma-separated list of component IDs (e.g. machine_id from Carbide)

Component types (required, currently only compute is supported):
  --type compute     : Compute nodes

Output formats:
  --output json      : JSON format (default)
  --output table     : Table format

Examples:
  # Get actual compute components from racks by name
  rla component actual --rack-names "rack-1,rack-2" --type compute

  # Get actual components from rack by ID
  rla component actual --rack-ids "uuid-1,uuid-2" --type compute

  # Get actual components by component IDs
  rla component actual --component-ids "machine-1,machine-2" --type compute

  # Output as table
  rla component actual --rack-names "rack-1" --type compute --output table
`,
		Run: func(cmd *cobra.Command, args []string) {
			doGetActualComponents()
		},
	}

	actualRackIDs       string
	actualRackNames     string
	actualComponentIDs  string
	actualComponentType string
	actualOutput        string
	actualHost          string
	actualPort          int
)

func init() {
	componentCmd.AddCommand(actualCmd)

	actualCmd.Flags().StringVar(&actualRackIDs, "rack-ids", "", "Comma-separated list of rack UUIDs")
	actualCmd.Flags().StringVar(&actualRackNames, "rack-names", "", "Comma-separated list of rack names")
	actualCmd.Flags().StringVar(&actualComponentIDs, "component-ids", "", "Comma-separated list of component IDs")
	actualCmd.Flags().StringVarP(&actualComponentType, "type", "t", "", "Component type (required): compute")
	actualCmd.Flags().StringVarP(&actualOutput, "output", "o", "json", "Output format: json, table")
	actualCmd.Flags().StringVar(&actualHost, "host", "localhost", "RLA server host")
	actualCmd.Flags().IntVar(&actualPort, "port", 50051, "RLA server port")
}

func doGetActualComponents() {
	// Validate input - exactly one of rack-ids, rack-names, or component-ids must be provided
	optionCount := 0
	if actualRackIDs != "" {
		optionCount++
	}
	if actualRackNames != "" {
		optionCount++
	}
	if actualComponentIDs != "" {
		optionCount++
	}

	if optionCount == 0 {
		log.Fatal().Msg("One of --rack-ids, --rack-names, or --component-ids must be specified")
	}
	if optionCount > 1 {
		log.Fatal().Msg("Only one of --rack-ids, --rack-names, or --component-ids can be specified")
	}

	// Component type is required for actual components
	if actualComponentType == "" {
		log.Fatal().Msg("--type is required (currently only 'compute' is supported)")
	}

	// Parse and validate component type
	componentType := parseComponentTypeToTypes(actualComponentType)
	if componentType != types.ComponentTypeCompute {
		log.Fatal().Msg("Only 'compute' component type is supported for actual components")
	}

	// Create client
	c, err := client.New(client.Config{
		Host: actualHost,
		Port: actualPort,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create client")
	}
	defer c.Close()

	ctx := context.Background()
	var result *client.GetActualComponentsResult

	// Call the appropriate client method based on input
	if actualRackIDs != "" {
		rackIDStrs := strings.Split(actualRackIDs, ",")
		rackIDs := make([]uuid.UUID, 0, len(rackIDStrs))
		for _, idStr := range rackIDStrs {
			id, err := uuid.Parse(strings.TrimSpace(idStr))
			if err != nil {
				log.Fatal().Err(err).Str("id", idStr).Msg("Invalid rack UUID")
			}
			rackIDs = append(rackIDs, id)
		}
		result, err = c.GetActualComponentsByRackIDs(ctx, rackIDs, componentType)
	} else if actualRackNames != "" {
		rackNames := strings.Split(actualRackNames, ",")
		for i := range rackNames {
			rackNames[i] = strings.TrimSpace(rackNames[i])
		}
		result, err = c.GetActualComponentsByRackNames(ctx, rackNames, componentType)
	} else if actualComponentIDs != "" {
		componentIDs := strings.Split(actualComponentIDs, ",")
		for i := range componentIDs {
			componentIDs[i] = strings.TrimSpace(componentIDs[i])
		}
		result, err = c.GetActualComponentsByComponentIDs(ctx, componentIDs, componentType)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get actual components")
	}

	// Output results
	switch actualOutput {
	case "json":
		outputActualJSON(result)
	case "table":
		outputActualTable(result)
	default:
		log.Fatal().Str("format", actualOutput).Msg("Unknown output format")
	}
}

func outputActualJSON(result *client.GetActualComponentsResult) {
	output := struct {
		Total      int         `json:"total"`
		Components interface{} `json:"components"`
	}{
		Total:      result.Total,
		Components: result.Components,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal JSON")
	}
	fmt.Println(string(data))
}

func outputActualTable(result *client.GetActualComponentsResult) {
	fmt.Printf("Total: %d components\n", result.Total)
	fmt.Println(strings.Repeat("-", 130))
	fmt.Printf("%-30s %-12s %-15s %-10s %-12s %-20s\n",
		"COMPONENT_ID", "TYPE", "FIRMWARE", "POWER", "HEALTH", "LAST_SEEN")
	fmt.Println(strings.Repeat("-", 130))

	for _, comp := range result.Components {
		lastSeen := ""
		if !comp.LastSeen.IsZero() {
			lastSeen = comp.LastSeen.Format("2006-01-02 15:04")
		}

		var compType string
		switch comp.Type {
		case types.ComponentTypeNVSwitch:
			compType = "nvlswitch"
		case types.ComponentTypePowerShelf:
			compType = "powershelf"
		default:
			compType = "compute"
		}

		fmt.Printf("%-30s %-12s %-15s %-10s %-12s %-20s\n",
			comp.ComponentID,
			compType,
			comp.FirmwareVersion,
			comp.PowerState,
			comp.HealthStatus,
			lastSeen,
		)
	}
}

// parseComponentTypeToTypes converts string to types.ComponentType
func parseComponentTypeToTypes(s string) types.ComponentType {
	switch strings.ToLower(s) {
	case "compute":
		return types.ComponentTypeCompute
	case "nvlswitch", "nvl-switch":
		return types.ComponentTypeNVSwitch
	case "powershelf", "power-shelf":
		return types.ComponentTypePowerShelf
	case "torswitch", "tor-switch":
		return types.ComponentTypeTORSwitch
	case "ums":
		return types.ComponentTypeUMS
	case "cdu":
		return types.ComponentTypeCDU
	default:
		return types.ComponentTypeUnknown
	}
}
