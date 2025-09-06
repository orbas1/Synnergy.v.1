package cli

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy"
)

var geoNode = synnergy.NewGeospatialNode()

func init() {
	cmd := &cobra.Command{
		Use:   "geospatial",
		Short: "Geospatial node operations",
	}

	recordCmd := &cobra.Command{
		Use:   "record [subject] [lat] [lon]",
		Args:  cobra.ExactArgs(3),
		Short: "Record a geospatial point",
		Run: func(cmd *cobra.Command, args []string) {
			lat, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				printOutput("invalid latitude: " + err.Error())
				return
			}
			lon, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				printOutput("invalid longitude: " + err.Error())
				return
			}
			gasPrint("GeospatialRecord")
			geoNode.Record(args[0], lat, lon)
			printOutput("recorded")
		},
	}

	historyCmd := &cobra.Command{
		Use:   "history [subject]",
		Args:  cobra.ExactArgs(1),
		Short: "Show recorded locations",
		Run: func(cmd *cobra.Command, args []string) {
			recs := geoNode.History(args[0])
			out := make([]map[string]any, len(recs))
			for i, r := range recs {
				out[i] = map[string]any{
					"subject":   r.Subject,
					"lat":       r.Latitude,
					"lon":       r.Longitude,
					"timestamp": r.Timestamp.Format(time.RFC3339),
				}
			}
			gasPrint("GeospatialHistory")
			printOutput(out)
		},
	}

	cmd.AddCommand(recordCmd, historyCmd)
	rootCmd.AddCommand(cmd)
}
