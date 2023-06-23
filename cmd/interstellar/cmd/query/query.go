package querycli

import (
	"fmt"

	"github.com/doggystylez/interstellar/client/cosmos/query"
	"github.com/doggystylez/interstellar/cmd/interstellar/cmd/flags"
	"github.com/doggystylez/interstellar/types"
	"github.com/spf13/cobra"
)

func QueryCmd() (qyCmd *cobra.Command) {
	qyCmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Query chain via gRPC",
		Long:    "Query chain via gRPC",
	}
	cmds := flags.AddFlags([]*cobra.Command{accountCmd(), addressCmd(), balanceCmd()}, flags.KeyFlags, flags.GlobalFlags)
	cmds = flags.AddFlags(append(cmds, chainCmd()), flags.QueryFlags)
	qyCmd.AddCommand(cmds...)
	return
}

func chainCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "chain-id",
		Short: "Query chain-id",
		Long:  "Query chain-id",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			rpc, err := flags.ProcessQueryFlags(cmd)
			if err != nil {
				panic(err)
			}
			chainId, err := query.GetChainId(rpc)
			if err != nil {
				panic(err)
			}
			fmt.Println(query.Jsonify(chainId))
		},
	}
	return
}

func accountCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "account <address>",
		Short: "Query account info",
		Long:  "Query account info by address, keyname, or privkey",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := flags.ProcessGlobalFlags(cmd)
			if err != nil {
				panic(err)
			}
			config.Rpc, err = flags.ProcessQueryFlags(cmd)
			if err != nil {
				panic(err)
			}
			if len(args) == 1 {
				config.TxInfo.Address = args[0]
			} else {
				config.TxInfo.KeyInfo.KeyRing, err = flags.ProcessKeyFlags(cmd)
				if err != nil {
					panic(err)
				}
				err = flags.CheckAddress(&config)
				if err != nil {
					panic(err)
				}
			}
			account, err := query.GetAccountInfoFromAddress(config)
			if err != nil {
				panic(err)
			}
			fmt.Println(query.Jsonify(account))
		},
	}
	return
}

func addressCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "address",
		Short: "Query account address",
		Long:  "Query account address by keyname or privkey",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := flags.ProcessGlobalFlags(cmd)
			if err != nil {
				panic(err)
			}
			config.Rpc, err = flags.ProcessQueryFlags(cmd)
			if err != nil {
				panic(err)
			}
			config.TxInfo.KeyInfo.KeyRing, err = flags.ProcessKeyFlags(cmd)
			if err != nil {
				panic(err)
			}
			err = flags.CheckAddress(&config)
			if err != nil {
				panic(err)
			}
			fmt.Println(query.Jsonify(types.AddressRes{Address: config.TxInfo.Address}))
		},
	}
	return
}

func balanceCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "balance <address>",
		Short: "Query account balance",
		Long:  "Query account balance, with optional denom filter",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				balances interface{}
			)
			config, err := flags.ProcessGlobalFlags(cmd)
			if err != nil {
				panic(err)
			}
			config.Rpc, err = flags.ProcessQueryFlags(cmd)
			if err != nil {
				panic(err)
			}
			if len(args) == 1 {
				config.TxInfo.Address = args[0]
			} else {
				config.TxInfo.KeyInfo.KeyRing, err = flags.ProcessKeyFlags(cmd)
				if err != nil {
					panic(err)
				}
				err = flags.CheckAddress(&config)
				if err != nil {
					panic(err)
				}
			}
			denom, err := cmd.Flags().GetString("denom")
			if err != nil {
				return
			}
			if denom == "" {
				resp, err := query.GetAllBalances(config.TxInfo.Address, config.Rpc)
				if err != nil {
					panic(err)
				}
				balances = resp.Balances

			} else {
				resp, err := query.GetBalanceByDenom(config.TxInfo.Address, denom, config.Rpc)
				if err != nil {
					panic(err)
				}
				balances = resp.Balance
			}
			fmt.Println(query.Jsonify(balances))
		},
	}
	cmd.Flags().StringP("denom", "d", "", "denom")
	return
}
