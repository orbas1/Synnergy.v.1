package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	syn "synnergy"
)

var geoNode = syn.NewGeospatialNode()

func init() {
	cmd := &cobra.Command{
		Use:   "geospatial",
		Short: "Record and query geospatial data",
	}

	recordCmd := &cobra.Command{
		Use:   "record [subject] [lat] [lon]",
		Args:  cobra.ExactArgs(3),
		Short: "Record a geospatial data point",
		RunE: func(cmd *cobra.Command, args []string) error {
			lat, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			lon, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return err
			}
			geoNode.Record(args[0], lat, lon)
			fmt.Println("recorded")
			return nil
		},
	}

	historyCmd := &cobra.Command{
		Use:   "history [subject]",
		Args:  cobra.ExactArgs(1),
		Short: "Show recorded locations for a subject",
		Run: func(cmd *cobra.Command, args []string) {
			for _, r := range geoNode.History(args[0]) {
				fmt.Printf("%s %.6f %.6f %s\n", r.Subject, r.Latitude, r.Longitude, r.Timestamp.Format(time.RFC3339))
			}
		},
	}

	cmd.AddCommand(recordCmd, historyCmd)
	rootCmd.AddCommand(cmd)
}
