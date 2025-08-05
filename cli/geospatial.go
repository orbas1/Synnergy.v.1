package cli

import (
	"fmt"
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
				fmt.Println("invalid latitude:", err)
				return
			}
			lon, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Println("invalid longitude:", err)
				return
			}
			geoNode.Record(args[0], lat, lon)
		},
	}

	historyCmd := &cobra.Command{
		Use:   "history [subject]",
		Args:  cobra.ExactArgs(1),
		Short: "Show recorded locations",
		Run: func(cmd *cobra.Command, args []string) {
			recs := geoNode.History(args[0])
			for _, r := range recs {
				fmt.Printf("%s %f %f %s\n", r.Subject, r.Latitude, r.Longitude, r.Timestamp.Format(time.RFC3339))
			}
		},
	}

	cmd.AddCommand(recordCmd, historyCmd)
	rootCmd.AddCommand(cmd)
}
