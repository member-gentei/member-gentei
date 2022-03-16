/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// adhocCmd represents the adhoc command
var adhocCmd = &cobra.Command{
	Use:   "adhoc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx      = context.Background()
			db       = mustOpenDB(ctx)
			existing = map[string]bool{}
			toRemove []int
			afterID  int
		)
		for {
			ums := db.UserMembership.Query().
				Where(usermembership.IDGT(afterID)).
				WithUser().
				WithYoutubeTalent().
				Order(ent.Asc(usermembership.FieldID)).
				Limit(4000).
				AllX(ctx)
			for _, um := range ums {
				if um.Edges.User == nil {
					toRemove = append(toRemove, um.ID)
					continue
				}
				key := fmt.Sprintf("%d-%s", um.Edges.User.ID, um.Edges.YoutubeTalent.ID)
				if existing[key] {
					toRemove = append(toRemove, um.ID)
				} else {
					existing[key] = true
				}
			}
			if len(ums) < 4000 {
				break
			}
			afterID = ums[len(ums)-1].ID
			log.Info().Int("toRemoveCount", len(toRemove)).Int("after", afterID).Msg("adhoc")
		}
		log.Info().Int("toRemoveCount", len(toRemove)).Msg("adhoc: removing duplicates and SET NULL'd")
		for i := 0; i < len(toRemove); i += 1000 {
			var batch []int
			if i+1000 > len(toRemove) {
				batch = toRemove[i:]
			} else {
				batch = toRemove[i : i+1000]
			}
			db.UserMembership.Delete().
				Where(usermembership.IDIn(batch...)).
				ExecX(ctx)
		}
	},
}

func init() {
	adminCmd.AddCommand(adhocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adhocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adhocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
